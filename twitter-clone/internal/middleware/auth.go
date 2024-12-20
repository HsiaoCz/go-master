package middleware

import (
	"net/http"
	"strings"
	"twitter-clone/internal/utils"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}
