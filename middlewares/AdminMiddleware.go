package middleware

import (
	"github.com/labstack/echo/v4"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user's role from the JWT token.

		role := c.Get("role")

		// Convert the `role` variable to type `float64`.
		roleFloat := role.(float64)

		if roleFloat != 1 {
			c.Error(echo.ErrForbidden)
			return echo.ErrForbidden
		}

		// Otherwise, continue with the request.
		return next(c)
	}
}
