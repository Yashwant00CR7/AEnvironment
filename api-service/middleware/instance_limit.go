/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	backendmodels "envhub/models"
	"fmt"

	"api-service/service"
	"api-service/util"

	"github.com/gin-gonic/gin"
)

// InstanceLimitMiddleware authentication middleware based on Bearer Token
// Uses unified response format from models
func InstanceLimitMiddleware(client *service.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token not configured, skip validation
		token := util.GetCurrentToken(c)
		if token == nil {
			c.Next()
			return
		}

		// List function not configured, skip validation
		if client == nil {
			return
		}
		// Check if instance count exceeds limit
		count, err := client.GetCurrentInstanceCount(token.Token)
		if err != nil {
			backendmodels.JSONErrorWithMessage(c, 500, "Failed to get instance count: "+err.Error())
			c.Abort()
			return
		}
		// Exceeds limit
		if count >= token.MaxInstanceCount && token.MaxInstanceCount > 0 {
			backendmodels.JSONErrorWithMessage(c, 403, fmt.Sprintf("Instance count exceeds limit: %d", token.MaxInstanceCount))
			c.Abort()
			return
		}
		c.Next()
	}
}
