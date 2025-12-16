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
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenHandler interface {
	ValidateToken(token string) (*backendmodels.Token, error)
}

// AuthTokenMiddleware authentication middleware based on Bearer Token
// Uses unified response format from models
func AuthTokenMiddleware(tokenEnabled bool, tokenHandler TokenHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !tokenEnabled {
			defaultToken := backendmodels.GenerateToken("default", 0)
			c.Set("token", defaultToken)
			c.Next()
			return
		}
		// 1. Extract Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			backendmodels.JSONErrorWithMessage(c, 401, "Missing Authorization header")
			c.Abort()
			return
		}

		// 2. Check if it starts with Bearer
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			backendmodels.JSONErrorWithMessage(c, 401, "Authorization header format should be Bearer <token>")
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			backendmodels.JSONErrorWithMessage(c, 401, "Token does not exist")
			c.Abort()
			return
		}

		// 3. Call backendClient to validate token
		validateToken, err := tokenHandler.ValidateToken(token)
		if err != nil {
			backendmodels.JSONErrorWithMessage(c, 401, "Token validation failed: "+err.Error())
			c.Abort()
			return
		}

		if validateToken == nil {
			backendmodels.JSONErrorWithMessage(c, 401, "Invalid token")
			c.Abort()
			return
		}
		// 4. Validation passed, continue processing
		c.Set("token", validateToken)
		c.Next()
	}
}
