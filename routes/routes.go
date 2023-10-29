package routes

import (
	"ppdb_sekolah_go/constans"
	"ppdb_sekolah_go/controllers"
	m "ppdb_sekolah_go/middlewares"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	m.LogMiddleware(e)

	e.GET("/datapokok", controllers.GetDatapokokController)
	e.GET("/datapokok/:id", controllers.GetDatapokokControllerByID)
	e.POST("/datapokok", controllers.CreateDatapokokController)
	e.DELETE("/datapokok/:id", controllers.DeleteDatapokokController)
	e.PUT("/datapokok/:id", controllers.UpdateDatapokokController)

	e.GET("/config", controllers.GetConfigController)
	e.GET("/config/:id", controllers.GetConfigControllerByID)
	e.POST("/config", controllers.CreateConfigController)
	e.DELETE("/config/:id", controllers.DeleteConfigController)
	e.PUT("/config/:id", controllers.UpdateConfigController)

	e.GET("/nilai", controllers.GetNilaiController)
	e.GET("/nilai/:id", controllers.GetNilaiControllerById)
	e.POST("/nilai", controllers.CreateNilaiController)
	e.DELETE("/nilai/:id", controllers.DeleteNilaiController)
	e.PUT("/nilai/:id", controllers.UpdateNilaiController)

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

	return e
}
