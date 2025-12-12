// Copyright 2022 Woodpecker Authors
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

package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

func TestUsers(t *testing.T) {
	store, closer := newTestStore(t, new(model.Account), new(model.Org), new(model.Secret), new(model.Repo), new(model.Perm))
	defer closer()

	count, err := store.GetUserCount()
	assert.NoError(t, err)
	assert.Zero(t, count)

	user := model.Account{
		AccountName:   "joe",
		ForgeRemoteID: "joe",
		AccessToken:   "f0b461ca586c27872b43a0685cbc2847",
		RefreshToken:  "976f22a5eef7caacb7e678d6c52f49b1",
	}
	assert.NoError(t, err)

	err2 := store.UpdateUser(&user)
	assert.NoError(t, err2)

	getUser, err := store.GetUser(user.ID, false)
	assert.NoError(t, err)
	assert.Equal(t, user.AccountName, getUser.AccountName)
	assert.Equal(t, user.AccessToken, getUser.AccessToken)
	assert.Equal(t, user.RefreshToken, getUser.RefreshToken)

	getUser, err = store.GetUserByLogin(user.ForgeID, user.AccountName)
	assert.NoError(t, err)
	assert.Equal(t, user.AccountName, getUser.AccountName)

	// check unique login
	user2 := model.Account{
		AccountName:   "Joe",
		ForgeRemoteID: "joe",
		AccessToken:   "ab20g0ddaf012c744e136da16aa21ad9",
	}

	user2 = model.Account{
		AccountName:   "jane",
		ForgeRemoteID: "jane",
		AccessToken:   "ab20g0ddaf012c744e136da16aa21ad9",
		Hash:          "A",
	}
	users, err := store.GetUserList(&model.ListOptions{Page: 1, PerPage: 50})
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	// "jane" user is first due to alphabetic sorting
	assert.Equal(t, user2.AccountName, users[0].AccountName)
	assert.Equal(t, user2.AccessToken, users[0].AccessToken)

	count, err = store.GetUserCount()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, count)

	getUser, err1 := store.GetUser(user.ID, false)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	_, err3 := store.GetUser(getUser.ID, false)
	assert.Error(t, err3)
}

func TestCreateUserWithExistingOrg(t *testing.T) {
	store, closer := newTestStore(t, new(model.Account), new(model.Org), new(model.Perm))
	defer closer()

	existingOrg := &model.Org{
		ForgeID: 1,
		IsUser:  true,
		Name:    "existingOrg",
		Private: false,
	}

	err := store.OrgCreate(existingOrg)
	assert.NoError(t, err)
	assert.EqualValues(t, "existingOrg", existingOrg.Name)

	// Create a new user with the same name as the existing organization

	updatedOrg, err := store.OrgGet(existingOrg.ID)
	assert.NoError(t, err)
	assert.Equal(t, "existingOrg", updatedOrg.Name)

	newOrg, err := store.OrgFindByName("new-user", 1)
	assert.NoError(t, err)
	assert.Equal(t, "new-user", newOrg.Name)
}
