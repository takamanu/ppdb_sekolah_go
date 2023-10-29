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

func GetNilaiController(c echo.Context) error {
	var nilais []models.Nilai
	if err := configs.DB.Find(&nilais).Error; err != nil {
		log.Errorf("Failed to get nilai: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"nilais":  nilais,
	})
}

func GetNilaiControllerById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	var nilai models.Nilai
	if err := configs.DB.First(&nilai, id).Error; err != nil {
		log.Errorf("Failed to get nilai with id %d: %s", id, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"nilai":   nilai,
	})
}

func CreateNilaiController(c echo.Context) error {
	nilai := models.Nilai{}
	if err := c.Bind(&nilai); err != nil {
		log.Errorf("Failed` to bind request: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	loger.Println(nilai)

	if err := configs.DB.Create(&nilai).Error; err != nil {
		log.Errorf("Failed to create nilai: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new nilai",
		"nilai":   nilai,
	})
}

// delete nilai by id
func DeleteNilaiController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var nilai models.Nilai
	if err := configs.DB.First(&nilai, id).Error; err != nil {
		log.Errorf("Failed to get nilai with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusNotFound, "Nilai not found")
	}

	if err := configs.DB.Delete(&nilai).Error; err != nil {
		log.Errorf("Failed to delete nilai with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete nilai")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success deleted nilai",
	})
}

// update nilai by id
func UpdateNilaiController(c echo.Context) error {
	// get nilai id from url param
	nilaiId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid Nilai id")
	}

	// get nilai by id
	var nilai models.Nilai
	if err := configs.DB.First(&nilai, nilaiId).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Nilai not found")
	}

	// bind request body to nilai struct
	if err := c.Bind(&nilai); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// save nilai to database
	if err := configs.DB.Save(&nilai).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success updated",
		"nilai":   nilai,
	})
}
