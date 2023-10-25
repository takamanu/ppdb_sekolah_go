package routes

import (
	"ppdb_sekolah_go/controllers"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	eAuthBasic := e.Group("/auth")
	eAuthBasic.Use(mid.BasicAuth(m.BasicAuthDB))
	eAuthBasic.GET("/users", controllers.GetUsersController)

	eJwt := e.Group("/jwt")
	eJwt.Use(mid.JWT([]byte(constans.SECRET_JWT)))
	e.GET("/users", controllers.GetUsersController)
	e.GET("/users/:id", controllers.GetUserController)
	e.POST("/users", controllers.CreateUserController)
	e.DELETE("/users/:id", controllers.DeleteUserController)
	e.PUT("/users/:id", controllers.UpdateUserController)
}
