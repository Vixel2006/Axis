package middlewares

import (
	"net/http"
	"strings"

	"axis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func JWTAuth(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.With().Str("middleware", "JWTAuth").Logger()

		h := c.GetHeader("Authorization")
		if h == "" {
			log.Warn().Msg("Missing authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(h, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn().Str("header", h).Msg("Invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		claims, err := utils.ParseToken(parts[1])

		if err != nil {
			log.Warn().Err(err).Msg("Invalid or expired token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		log.Debug().Int("user_id", claims.UserID).Msg("User ID extracted from token")
		c.Set("user_id", claims.UserID)

		log.Debug().Msg("JWT authentication successful, continuing to next handler")
		c.Next()
	}
}
