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

package util

import (
	"envhub/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// SplitEnvNameVersion splits env string by @ into name and version
// Example: "leopard-linux-v1@1.0.0" -> name: "leopard-linux-v1", version: "1.0.0"
// If there's no @ symbol, the entire string is used as name, version is empty
func SplitEnvNameVersion(env string) (name, version string) {
	if env == "" {
		return "", ""
	}

	parts := strings.Split(env, "@")
	if len(parts) == 1 {
		// No @ symbol, use entire string as name
		return parts[0], ""
	} else if len(parts) >= 2 {
		// Has @ symbol, first part as name, second part as version
		// If there are multiple @ symbols, only take the first two parts
		return parts[0], parts[1]
	}

	return env, ""
}

// SplitEnvNameVersionStrict splits env string in strict mode
// Requires @ symbol to be present, otherwise returns error
func SplitEnvNameVersionStrict(env string) (name, version string, err error) {
	if env == "" {
		return "", "", fmt.Errorf("env string is empty")
	}

	parts := strings.Split(env, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid env format, expected 'name@version', got: %s", env)
	}

	if parts[0] == "" {
		return "", "", fmt.Errorf("env name cannot be empty")
	}

	if parts[1] == "" {
		return "", "", fmt.Errorf("env version cannot be empty")
	}

	return parts[0], parts[1], nil
}

// JoinEnvNameVersion joins name and version with @
func JoinEnvNameVersion(name, version string) string {
	if version == "" {
		return name
	}
	return fmt.Sprintf("%s@%s", name, version)
}

func GetCurrentToken(c *gin.Context) *models.Token {
	token, exists := c.Get("token")
	if !exists {
		return nil
	}
	return token.(*models.Token)
}
