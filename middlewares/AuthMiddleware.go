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
			return jsonResponse(c, http.StatusUnauthorized, false, "Invalid token", nil)

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
			return jsonResponse(c, http.StatusUnauthorized, false, err.Error(), nil)
		}

		// Access the claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)

		// Check if the JWT token claims are valid

		if !ok {
			return jsonResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		}

		// Retrieve the userId from the claims
		userId, ok := claims["userId"].(float64)

		// Check if the userId is found in the claims

		if !ok {
			return jsonResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		}

		role, ok := claims["role"].(float64)

		// Check if the userId is found in the claims

		if !ok {
			return jsonResponse(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		}

		// Set the userId in the context
		c.Set("userId", userId)
		c.Set("role", role)

		// Call the next handler

		return next(c)
	}
}

func jsonResponse(c echo.Context, status int, success bool, message string, data interface{}) error {
	return c.JSON(status, map[string]interface{}{
		constans.SUCCESS: success,
		constans.MESSAGE: message,
		constans.DATA:    data,
	})
}
