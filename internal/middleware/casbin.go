package middleware

import (
	"net/http"

	"simple-go/pkg/logger"
	"simple-go/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Context keys for resource and action
const (
	ResourceKey = "casbin_resource"
	ActionKey   = "casbin_action"
)

// SetReSource404NovelDownloaderction sets the resource and action in the context for Casbin authorization
// This should be called in handlers before the Casbin middleware runs
func SetReSource404NovelDownloaderction(c *gin.Context, resource, action string) {
	c.Set(ResourceKey, resource)
	c.Set(ActionKey, action)
}

// GetReSource404NovelDownloaderction retrieves the resource and action from the context
func GetReSource404NovelDownloaderction(c *gin.Context) (resource, action string, exists bool) {
	resourceVal, resourceExists := c.Get(ResourceKey)
	actionVal, actionExists := c.Get(ActionKey)

	if !resourceExists || !actionExists {
		return "", "", false
	}

	resource, ok1 := resourceVal.(string)
	action, ok2 := actionVal.(string)

	if !ok1 || !ok2 {
		return "", "", false
	}

	return resource, action, true
}

// CasbinAuthorizer creates a middleware that checks permissions using resource-action approach
func CasbinAuthorizer(enforcer *casbin.Enforcer, roleGetter func(*gin.Context) ([]string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := roleGetter(c)
		if err != nil {
			logger.Error(err, "failed to get user roles")
			response.Error(c, http.StatusInternalServerError, "Failed to get user roles")
			c.Abort()
			return
		}

		if len(roles) == 0 {
			response.Error(c, http.StatusForbidden, "User has no roles assigned")
			c.Abort()
			return
		}

		// Get resource and action from context
		resource, action, exists := GetReSource404NovelDownloaderction(c)
		if !exists {
			logger.Error(nil, "resource and action not set in context")
			response.Error(c, http.StatusInternalServerError, "Authorization configuration error")
			c.Abort()
			return
		}

		// Check if any of the user's roles have permission
		allowed := false
		for _, role := range roles {
			ok, err := enforcer.Enforce(role, resource, action)
			if err != nil {
				logger.Error(err, "authorization check failed")
				response.Error(c, http.StatusInternalServerError, "Authorization check failed")
				c.Abort()
				return
			}
			if ok {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, http.StatusForbidden, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission is a convenience middleware that sets the resource and action, then checks authorization
func RequirePermission(resource, action string, enforcer *casbin.Enforcer, roleGetter func(*gin.Context) ([]string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetReSource404NovelDownloaderction(c, resource, action)
		CasbinAuthorizer(enforcer, roleGetter)(c)
	}
}
