package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/user"
	"google.golang.org/grpc/status"
)

const (
	UserDataKey       = "User"
	Authorization     = "Authorization"
	Bearer            = "Bearer"
	UserProfileCtxKey = "UserProfile"
)

func JwtMiddleware(userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := extractTokenFromRequestHeader(c.Request)
		if jwtToken == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"success": false,
					"error": gin.H{
						"code":    http.StatusUnauthorized,
						"message": "token is empty",
					},
				})
			return
		}
		userData, err := userService.GetUserData(c, jwtToken, c.FullPath())
		if err != nil {
			e, ok := status.FromError(err)
			if !ok {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					gin.H{
						"success": false,
						"error": gin.H{
							"code":    http.StatusForbidden,
							"message": "cannot get user detail",
						},
					})
				return
			}
			c.AbortWithStatusJSON(
				int(e.Code()),
				gin.H{
					"success": false,
					"error": gin.H{
						"code":    int(e.Code()),
						"message": "cannot get user detail",
					},
				})
			return
		}
		c.Set(UserDataKey, userData)
		c.Next()
	}
}

func extractTokenFromRequestHeader(r *http.Request) string {
	token := r.Header.Get(Authorization)
	token = strings.Trim(token, `"`)
	if strings.Contains(token, Bearer) {
		token = strings.TrimPrefix(token, Bearer+" ")
	}
	return token
}
