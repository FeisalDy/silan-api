package middleware

import (
	"net/http"
	"strings"

	"simple-go/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinAuthorizer is a middleware that enforces Casbin policies
func CasbinAuthorizer(enforcer *casbin.Enforcer, roleGetter func(*gin.Context) ([]string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user roles
		roles, err := roleGetter(c)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to get user roles", err)
			c.Abort()
			return
		}

		if len(roles) == 0 {
			response.Error(c, http.StatusForbidden, "User has no roles assigned", nil)
			c.Abort()
			return
		}

		// Get request path and method
		path := c.Request.URL.Path
		method := c.Request.Method

		// Normalize path to match Casbin policy pattern (replace dynamic segments)
		normalizedPath := normalizePath(path)

		// Check if any of the user's roles have permission
		allowed := false
		for _, role := range roles {
			ok, err := enforcer.Enforce(role, normalizedPath, method)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, "Authorization check failed", err)
				c.Abort()
				return
			}
			if ok {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, http.StatusForbidden, "Access denied", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// normalizePath converts actual paths to policy patterns
// E.g., /api/v1/users/123 -> /api/v1/users/:id
func normalizePath(path string) string {
	parts := strings.Split(path, "/")

	// Handle /api/v1/users/:id pattern
	if len(parts) >= 4 && parts[1] == "api" && parts[2] == "v1" {
		resource := parts[3]

		// If there's a 4th segment, it's likely an ID
		if len(parts) >= 5 && parts[4] != "" {
			return "/api/v1/" + resource + "/:id"
		}
	}

	return path
}
