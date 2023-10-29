package controllers

import (
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/models"

	loger "log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GetConfigController(c echo.Context) error {
	var users []models.Config
	if err := configs.DB.Find(&users).Error; err != nil {
		log.Errorf("Failed to get config: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"users":   users,
	})
}

func GetConfigControllerByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	var user models.Config
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get config with id %d: %s", id, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"user":    user,
	})
}

func CreateConfigController(c echo.Context) error {
	user := models.Config{}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Failed to bind request: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	loger.Println(user)

	if err := configs.DB.Create(&user).Error; err != nil {
		log.Errorf("Failed to create config: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new config",
		"user":    user,
	})
}

// delete user by id
func DeleteConfigController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var user models.Datapokok
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get config with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusNotFound, "Config not found")
	}

	if err := configs.DB.Delete(&user).Error; err != nil {
		log.Errorf("Failed to delete config with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete config")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success deleted config",
	})
}

// update user by id
func UpdateConfigController(c echo.Context) error {
	// get user id from url param
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid config id")
	}

	// get user by id
	var user models.Datapokok
	if err := configs.DB.First(&user, userId).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Config not found")
	}

	// bind request body to user struct
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// save user to database
	if err := configs.DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success updated",
		"user":    user,
	})
}
