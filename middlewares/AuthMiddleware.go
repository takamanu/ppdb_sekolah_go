package middleware

import (
	"fmt"
	"net/http"
	"ppdb_sekolah_go/constans"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Get the authorization header
		authorizationHeader := c.Request().Header.Get("Authorization")

		// Check if the authorization header is valid
		if !strings.Contains(authorizationHeader, "Bearer") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Extract the JWT token from the authorization header
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Signing method invalid")
			}

			return []byte(constans.SECRET_JWT), nil
		})

		// Check if the JWT token is valid

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		// Access the claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)

		// Check if the JWT token claims are valid

		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "JWT claims not found")
		}

		// Retrieve the userId from the claims
		userId, ok := claims["userId"].(float64)

		// Check if the userId is found in the claims

		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "userId not found in JWT claims")
		}

		role, ok := claims["role"].(float64)

		// Check if the userId is found in the claims

		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Role not found in JWT claims")
		}

		// Set the userId in the context
		c.Set("userId", userId)
		c.Set("role", role)

		// Call the next handler

		return next(c)
	}
}
