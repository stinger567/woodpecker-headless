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
	"net/http"

	"github.com/gin-gonic/gin"

	"go.woodpecker-ci.org/woodpecker/v3/server/store"
)

// Health
//
//	@Summary		Health information
//	@Description	If everything is fine, just a 204 will be returned, a 500 signals server state is unhealthy.
//	@Router			/healthz [get]
//	@Produce		plain
//	@Success		204
//	@Failure		500
//	@Tags			System
func Health(c *gin.Context) {
	if err := store.FromContext(c).Ping(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// Version
//
//	@Summary		Get version
//	@Description	Endpoint returns the server version and build information.
//	@Router			/version [get]
//	@Produce		json
//	@Success		200	{object}	object{source=string,version=string}
//	@Tags			System
func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"source":  "https://github.com/woodpecker-ci/woodpecker",
		"version": "0.0.0",
	})
}
