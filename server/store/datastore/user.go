// Copyright 2021 Woodpecker Authors
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
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
)

func (s storage) GetUser(id string, internal bool) (*model.Account, error) {
	user := new(model.Account)
	return user, wrapGet(s.engine.ID(id).Get(user))
}

func (s storage) GetUserByRemoteID(forgeID string, userRemoteID model.ForgeRemoteID, internal bool) (*model.Account, error) {
	sess := s.engine.NewSession()
	user := new(model.Account)
	return user, wrapGet(sess.Where("forge_id = ? AND forge_remote_id = ?", forgeID, userRemoteID).Get(user))
}

func (s storage) GetUserByLogin(forgeID string, login string, internal bool) (*model.Account, error) {
	sess := s.engine.NewSession()
	user := new(model.Account)
	return user, wrapGet(sess.Where("forge_id = ? AND login=?", forgeID, login).Get(user))
}

func (s storage) GetUserList(p *model.ListOptions) ([]*model.Account, error) {
	var users []*model.Account
	return users, s.paginate(p).OrderBy("login").Find(&users)
}

func (s storage) GetUserCount() (int64, error) {
	return s.engine.Count(new(model.Account))
}

func (s storage) UpdateUser(user *model.Account) error {
	_, err := s.engine.ID(user.ForgeRemoteID).AllCols().Update(user)
	return err
}
