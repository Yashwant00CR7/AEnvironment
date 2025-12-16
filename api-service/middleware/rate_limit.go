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
	"envhub/models"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RateLimit handles QPS rate limiting for apiserver requests
// If qps is 0, no rate limiting is performed; otherwise rate limiting is executed
//
// If the rate limit total is exceeded within the expected time, returns 429 Too Many Requests HTTP status code,
// and adds rate limit related headers to response header:
// - X-Rate-Limit-Limit: rate limit value (maximum number of requests within specified time)
// - RateLimit-Remaining: remaining number of requests that can be made
// - RateLimit-Reset: time point when rate limit resets, value is unix epoch timestamp
//
// Current rate limiting implementation only considers total QPS of a single faas-api-service,
// more detailed rate limiting strategies can be added later:
// 1. Distributed QPS total limit
// 2. Different rate limiting strategies based on user, tenant, etc.
// 3. Different priority rate limiting strategies based on different API interfaces
func RateLimit(qps int64) gin.HandlerFunc {
	log.Infof("initial rate limiter with qps: %d", qps)
	var lmt *limiter.Limiter = nil
	if qps > 0 {
		lmt = tollbooth.NewLimiter(float64(qps), nil)
	}

	return func(c *gin.Context) {
		if lmt == nil {
			c.Next()
			return
		}

		httpError := tollbooth.LimitByKeys(lmt, []string{"total_qps"})
		if httpError != nil {
			response := models.NewErrorResponseWithData(httpError.StatusCode, httpError.Message)
			c.JSON(httpError.StatusCode, response)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
