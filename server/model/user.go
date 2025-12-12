// Copyright 2021 Woodpecker Authors
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

package model

import (
	"errors"
	"regexp"
)

// Validate a username (e.g. from github).
var reUsername = regexp.MustCompile("^[a-zA-Z0-9-_.]+$")

var errUserLoginInvalid = errors.New("invalid user login")

const maxLoginLen = 250

// Account represents a registered user.
type Account struct {
	ID      string
	ForgeID string `json:"forge_id,omitempty" xorm:"forge_id UNIQUE(forge_id) UNIQUE(forge_login)"`

	ForgeRemoteID ForgeRemoteID `json:"forge_remote_id" xorm:"forge_remote_id UNIQUE(forge_id)"`

	// AccountName is the username for this user.
	AccountName string `json:"login"  xorm:"'login' UNIQUE(forge_login)"`

	// AccessToken is the oauth2 access token.
	AccessToken string `json:"-"  xorm:"TEXT 'access_token'"`

	// RefreshToken is the oauth2 refresh token.
	RefreshToken string `json:"-" xorm:"TEXT 'refresh_token'"`

	// Expiry is the AccessToken expiration timestamp (unix seconds).
	Expiry int64 `json:"-" xorm:"expiry"`

	// Hash is a unique token used to sign tokens.
	Hash string `json:"-" xorm:"UNIQUE varchar(500) 'hash'"`

	// OrgID is the of the user as model.Org.
	OrgID string `json:"org_id" xorm:"org_id"`

	Internal bool
} //	@name	User

// TableName return database table name for xorm.
func (Account) TableName() string {
	return "users"
}

// Validate validates the required fields and formats.
func (account *Account) Validate() error {
	switch {
	case len(account.AccountName) == 0:
		return errUserLoginInvalid
	case len(account.AccountName) > maxLoginLen:
		return errUserLoginInvalid
	case !reUsername.MatchString(account.AccountName):
		return errUserLoginInvalid
	default:
		return nil
	}
}
