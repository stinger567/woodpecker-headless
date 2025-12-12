// Copyright 2018 Drone.IO Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/tink/go/subtle/random"
	"github.com/rs/zerolog/log"

	"go.woodpecker-ci.org/woodpecker/v3/server"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"go.woodpecker-ci.org/woodpecker/v3/server/store"
	"go.woodpecker-ci.org/woodpecker/v3/server/store/types"
	"go.woodpecker-ci.org/woodpecker/v3/shared/httputil"
	"go.woodpecker-ci.org/woodpecker/v3/shared/token"
	"go.woodpecker-ci.org/woodpecker/v3/shared/utils"
)

const (
	stateTokenDuration = time.Minute * 5
	perPage            = 50
	maxPage            = 10000
)

func HandleAuth(c *gin.Context) {
	// TODO: check if this is really needed
	c.Writer.Header().Del("Content-Type")

	// redirect when getting oauth error from forge to login page
	if err := c.Request.FormValue("error"); err != "" {
		query := url.Values{}
		query.Set("error", err)
		if errorDescription := c.Request.FormValue("error_description"); errorDescription != "" {
			query.Set("error_description", errorDescription)
		}
		if errorURI := c.Request.FormValue("error_uri"); errorURI != "" {
			query.Set("error_uri", errorURI)
		}
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/login?%s", server.Config.Server.RootPath, query.Encode()))
		return
	}

	_store := store.FromContext(c)

	code := c.Request.FormValue("code")
	state := c.Request.FormValue("state")
	isCallback := code != "" && state != ""
	var forgeID string

	if isCallback { // validate the state token
		stateToken, err := token.Parse([]token.Type{token.OAuthStateToken}, state, func(_ *token.Token) (string, error) {
			return server.Config.Server.JWTSecret, nil
		})
		if err != nil {
			log.Error().Err(err).Msg("cannot verify state token")
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=invalid_state")
			return
		}

		_forgeID := stateToken.Get("forge-id")
		forgeID = _forgeID
		if err != nil {
			log.Error().Err(err).Msg("forge-id of state token invalid")
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=invalid_state")
			return
		}
	} else { // only generate a state token if not a callback
		var err error

		_forgeID := c.Request.FormValue("forge_id")
		if _forgeID == "" {
			forgeID = "1" // fallback to main forge
		} else {
			forgeID = _forgeID
			if err != nil {
				log.Error().Err(err).Msg("forge-id of state token invalid")
				c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=invalid_state")
				return
			}
		}

		jwtSecret := server.Config.Server.JWTSecret
		exp := time.Now().Add(stateTokenDuration).Unix()
		stateToken := token.New(token.OAuthStateToken)
		stateToken.Set("forge-id", forgeID)
		state, err = stateToken.SignExpires(jwtSecret, exp)
		if err != nil {
			log.Error().Err(err).Msg("cannot create state token")
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
			return
		}
	}

	_forge, err := server.Config.Services.Manager.ForgeByID(forgeID, false)
	if err != nil {
		log.Error().Err(err).Msgf("cannot get forge by id %d", forgeID)
		c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
		return
	}

	userFromForge, redirectURL, err := _forge.Login(c, &forge_types.OAuthRequest{
		Code:  c.Request.FormValue("code"),
		State: state,
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot authenticate user")
		c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=oauth_error")
		return
	}
	// The user is not authorized yet -> redirect
	if userFromForge == nil {
		http.Redirect(c.Writer, c.Request, redirectURL, http.StatusSeeOther)
		return
	}

	// if organization filter is enabled, we need to check if the user is a member of one
	// of the configured organizations
	if server.Config.Permissions.Orgs.IsConfigured {
		isMember := false
		for page := 1; page <= maxPage; page++ {
			teams, terr := _forge.Teams(c, userFromForge, &model.ListOptions{
				Page:    page,
				PerPage: perPage,
			})
			if terr != nil {
				log.Error().Err(terr).Msgf("cannot verify team membership for %s", userFromForge.AccountName)
				c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
				return
			}
			if server.Config.Permissions.Orgs.IsMember(teams) {
				isMember = true
				break
			}
		}
		if !isMember {
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=org_access_denied")
			return
		}
	}

	var user *model.Account

	// get the user from the database
	user, err = _store.GetUserByRemoteID(forgeID, userFromForge.ForgeRemoteID, userFromForge.Internal)
	if err != nil && !errors.Is(err, types.RecordNotExist) {
		log.Error().Err(err).Msgf("cannot get user %s", userFromForge.AccountName)
		c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
		return
	}
	// update user login (in case forge supports renaming)
	if user != nil {
		user.AccountName = userFromForge.AccountName
	}

	// re-try with login name
	if user == nil || errors.Is(err, types.RecordNotExist) {
		user, err = _store.GetUserByLogin(forgeID, userFromForge.AccountName, userFromForge.Internal)
		if err != nil && !errors.Is(err, types.RecordNotExist) {
			log.Error().Err(err).Msgf("cannot get user %s", userFromForge.AccountName)
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
			return
		}
	}

	if user == nil || errors.Is(err, types.RecordNotExist) {
		// if self-registration is disabled we should return a not authorized error
		if !server.Config.Permissions.Open && !server.Config.Permissions.Admins.IsAdmin(userFromForge) {
			log.Error().Msgf("cannot register %s. registration closed", userFromForge.AccountName)
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=registration_closed")
			return
		}

		// create the user account
		user = &model.Account{
			ForgeID:       forgeID,
			ForgeRemoteID: userFromForge.ForgeRemoteID,
			AccountName:   userFromForge.AccountName,
			AccessToken:   userFromForge.AccessToken,
			RefreshToken:  userFromForge.RefreshToken,
			Expiry:        userFromForge.Expiry,
			Hash: base32.StdEncoding.EncodeToString(
				random.GetRandomBytes(32),
			),
		}

	}

	// create or set the user's organization if it isn't linked yet
	if user.OrgID == "" {
		// check if an org with the same name exists already and assign it to the user if it does
		org, err := _store.OrgFindByName(user.AccountName, forgeID, user.Internal)
		if err != nil && !errors.Is(err, types.RecordNotExist) {
			log.Error().Err(err).Msgf("cannot get org for user %s", user.AccountName)
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
			return
		}

		// if an org with the same name exists => assign org to the user
		if err == nil && org != nil {
			org.IsUser = true
			user.OrgID = org.ID

			if err := _store.OrgUpdate(org); err != nil {
				log.Error().Err(err).Msgf("cannot assign user %s to existing org %d", user.AccountName, org.ID)
			}
		}

		// if still no org with the same name exists => create a new org
		if user.OrgID == "" || errors.Is(err, types.RecordNotExist) {
			org := &model.Org{
				Name:    user.AccountName,
				IsUser:  true,
				Private: false,
				ForgeID: user.ForgeID,
			}
			if err := _store.OrgCreate(org); err != nil {
				log.Error().Err(err).Msgf("cannot create org for user %s", user.AccountName)
			}
			user.OrgID = org.ID
		}
	} else {
		// update org name if necessary
		org, err := _store.OrgGet(user.OrgID, user.Internal)
		if err != nil {
			log.Error().Err(err).Msgf("cannot get org %d for user %s", user.OrgID, user.AccountName)
			c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
			return
		}
		if org != nil && org.Name != user.AccountName {
			org.Name = user.AccountName
			if err := _store.OrgUpdate(org); err != nil {
				log.Error().Err(err).Msgf("cannot update org %d name to user name %s", org.ID, user.AccountName)
			}
		}
	}

	// update the user meta data and authorization data.
	user.AccessToken = userFromForge.AccessToken
	user.RefreshToken = userFromForge.RefreshToken
	user.ForgeID = forgeID
	user.ForgeRemoteID = userFromForge.ForgeRemoteID
	user.AccountName = userFromForge.AccountName

	if err := _store.UpdateUser(user); err != nil {
		log.Error().Err(err).Msgf("cannot update user %s", user.AccountName)
		c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
		return
	}

	exp := time.Now().Add(server.Config.Server.SessionExpires).Unix()
	_token := token.New(token.SessToken)
	_token.Set("user-id", user.ID)
	tokenString, err := _token.SignExpires(user.Hash, exp)
	if err != nil {
		log.Error().Msgf("cannot create token for user %s", user.AccountName)
		c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
		return
	}

	err = updateRepoPermissions(c, user, _store, _forge)
	if err != nil {
		log.Error().Err(err).Msgf("cannot update repo permissions for user %s", user.AccountName)
		c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/login?error=internal_error")
		return
	}

	httputil.SetCookie(c.Writer, c.Request, "user_sess", tokenString)

	c.Redirect(http.StatusSeeOther, server.Config.Server.RootPath+"/")
}

func updateRepoPermissions(c *gin.Context, user *model.Account, _store store.Store, _forge forge.Forge) error {
	repos, err := utils.Paginate(func(page int) ([]*model.Repo, error) {
		return _forge.Repos(c, user, &model.ListOptions{
			Page:    page,
			PerPage: perPage,
		})
	}, maxPage)
	if err != nil {
		return err
	}

	for _, forgeRepo := range repos {
		dbRepo, err := _store.GetRepoForgeID(forgeRepo.ForgeRemoteID, forgeRepo.Internal)
		if err != nil && errors.Is(err, types.RecordNotExist) {
			continue
		}
		if err != nil {
			return err
		}

		if !dbRepo.IsActive {
			continue
		}

		log.Debug().Msgf("synced user permission for user %s and repo %s", user.AccountName, dbRepo.FullName)
		perm := forgeRepo.Perm
		perm.Repo = dbRepo
		perm.RepoID = dbRepo.ID
		perm.UserID = user.ID
		perm.Synced = time.Now().Unix()
		if err := _store.PermUpsert(perm); err != nil {
			return err
		}
	}

	return nil
}
