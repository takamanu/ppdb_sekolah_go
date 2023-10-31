package routes

import (
	"net/http"
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/constans"
	"ppdb_sekolah_go/controllers"
	m "ppdb_sekolah_go/middlewares"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	m.LogMiddleware(e)

	e.POST("/login", controllers.LoginUserController)
	e.POST("/register", controllers.CreateUserController)

	// eAuthBasic := e.Group("/auth")
	// eAuthBasic.Use(mid.BasicAuth(m.BasicAuthDB))
	// eAuthBasic.GET("/users", controllers.GetUsersController)

	eSiswa := e.Group("/siswa")
	eSiswa.Use(mid.JWT([]byte(constans.SECRET_JWT)))
	eSiswa.Use(m.AuthMiddleware)
	eSiswa.POST("/datapokok", CreateDatapokokHandlerSiswa)
	eSiswa.GET("/datapokok", controllers.GetDatapokokControllerSiswa)
	eSiswa.POST("/jurusan", controllers.AIController)
	// eSiswa.PUT("/datapokok", controllers.UpdateDatapokokController)

	eSiswa.GET("/pengumuman", controllers.GetPengumumanSiswa)

	eAdmin := e.Group("/admin")
	eAdmin.Use(mid.JWT([]byte(constans.SECRET_JWT)))
	eAdmin.Use(m.AuthMiddleware)
	eAdmin.Use(m.AdminMiddleware)
	eAdmin.GET("/users", controllers.GetUsersController)
	eAdmin.GET("/users/:id", controllers.GetUserController)
	eAdmin.POST("/users", controllers.CreateUserController)
	eAdmin.DELETE("/users/:id", controllers.DeleteUserController)
	eAdmin.PUT("/users/:id", controllers.UpdateUserController)

	eAdmin.GET("/nilai", controllers.GetNilaiController)
	eAdmin.GET("/nilai/:id", controllers.GetNilaiControllerById)
	eAdmin.POST("/nilai", controllers.CreateNilaiController)
	eAdmin.DELETE("/nilai/:id", controllers.DeleteNilaiController)
	eAdmin.PUT("/nilai/:id", controllers.UpdateNilaiController)

	// eAdmin.GET("/datapokok-jwt", controllers.GetDatapokokControllerByIDUseJWT)
	eAdmin.GET("/datapokok", controllers.GetDatapokokController)
	eAdmin.GET("/datapokok/:id", controllers.GetDatapokokControllerByID)
	eAdmin.POST("/datapokok", CreateDatapokokHandler)
	eAdmin.DELETE("/datapokok/:id", controllers.DeleteDatapokokController)
	eAdmin.PUT("/datapokok/:id", controllers.UpdateDatapokokController)

	eAdmin.GET("/config", controllers.GetConfigController)
	eAdmin.GET("/config/:id", controllers.GetConfigControllerByID)
	eAdmin.POST("/config", controllers.CreateConfigController)
	eAdmin.DELETE("/config/:id", controllers.DeleteConfigController)
	eAdmin.PUT("/config/:id", controllers.UpdateConfigController)

	return e
}

func CreateDatapokokHandler(c echo.Context) error {
	client, bucketName, err := configs.InitGCB()
	if err != nil {
		// Handle the error, e.g., return an internal server error
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to initialize Google Cloud Storage")
	}

	return controllers.CreateDatapokokController(c, client, bucketName)
}

func CreateDatapokokHandlerSiswa(c echo.Context) error {
	client, bucketName, err := configs.InitGCB()
	if err != nil {
		// Handle the error, e.g., return an internal server error
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to initialize Google Cloud Storage")
	}

	return controllers.CreateDatapokokControllerSiswa(c, client, bucketName)
}
