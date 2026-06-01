package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// BlockRoles aborta la petición si el rol del JWT está en la lista.
func BlockRoles(roles ...string) gin.HandlerFunc {
	blocked := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		blocked[r] = struct{}{}
	}
	return func(c *gin.Context) {
		roleVal, _ := c.Get("user_role")
		role, _ := roleVal.(string)
		if _, ok := blocked[role]; ok {
			utils.ErrorResponse(c, http.StatusForbidden, services.ErrForbidden.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
