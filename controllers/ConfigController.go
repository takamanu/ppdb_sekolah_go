package controllers

import (
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/constans"
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
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}
	return jsonResponse(c, http.StatusOK, true, "Success", users)
}

func GetConfigControllerByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return jsonResponse(c, http.StatusNotFound, true, "Invalid id", nil)
	}
	var user models.Config
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get config with id %d: %s", id, err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success",
		constans.DATA:    user,
	})
}

func CreateConfigController(c echo.Context) error {
	user := models.Config{}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Failed to bind request: %s", err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}

	loger.Println(user)

	if err := configs.DB.Create(&user).Error; err != nil {
		log.Errorf("Failed to create config: %s", err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success create new config",
		constans.DATA:    user,
	})
}

func DeleteConfigController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid id", nil)

	}

	var user models.Config
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get config with id %d: %v", id, err)
		return jsonResponse(c, http.StatusNotFound, false, "Config not found", nil)
	}

	if err := configs.DB.Delete(&user).Error; err != nil {
		log.Errorf("Failed to delete config with id %d: %v", id, err)
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete config", nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success deleted config",
	})
}

// update user by id
func UpdateConfigController(c echo.Context) error {
	// get user id from url param
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return jsonResponse(c, http.StatusNotFound, true, "invalid config id", nil)
	}

	// get user by id
	var user models.Config
	if err := configs.DB.First(&user, userId).Error; err != nil {
		return jsonResponse(c, http.StatusNotFound, true, "Config not found", nil)
	}

	// bind request body to user struct
	if err := c.Bind(&user); err != nil {
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}

	// save user to database
	if err := configs.DB.Save(&user).Error; err != nil {
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success updated",
		constans.DATA:    user,
	})
}
