package controllers

import (
	"fmt"
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/constans"
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
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success",
		constans.DATA:    nilais,
	})
}

func GetNilaiControllerById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid id", nil)
	}
	var nilai models.Nilai
	if err := configs.DB.First(&nilai, id).Error; err != nil {
		log.Errorf("Failed to get nilai with id %d: %s", id, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success",
		constans.DATA:    nilai,
	})
}

func CreateNilaiController(c echo.Context) error {
	nilai := models.Nilai{}
	if err := c.Bind(&nilai); err != nil {
		log.Errorf("Failed` to bind request: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	loger.Println(nilai)

	if err := configs.DB.Create(&nilai).Error; err != nil {
		log.Errorf("Failed to create nilai: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success create new nilai",
		constans.DATA:    nilai,
	})
}

// delete nilai by id
func DeleteNilaiController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid id", nil)
	}

	var nilai models.Nilai
	if err := configs.DB.First(&nilai, id).Error; err != nil {
		log.Errorf("Failed to get nilai with id %d: %v", id, err)
		return jsonResponse(c, http.StatusNotFound, false, "Nilai not found", nil)

	}

	if err := configs.DB.Delete(&nilai).Error; err != nil {
		log.Errorf("Failed to delete nilai with id %d: %v", id, err)
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete nilai", nil)

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success deleted nilai",
	})
}

// update nilai by id
func UpdateNilaiController(c echo.Context) error {
	// get nilai id from url param
	nilaiId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, "invalid Nilai id", nil)
	}

	// get nilai by id
	var nilai models.Nilai
	if err := configs.DB.First(&nilai, nilaiId).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, "Nilai not found", nil)
	}

	// bind request body to nilai struct
	if err := c.Bind(&nilai); err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	// save nilai to database
	if err := configs.DB.Save(&nilai).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success updated",
		constans.DATA:    nilai,
	})
}

func GetPengumumanSiswa(c echo.Context) error {

	userId := c.Get("userId")
	fmt.Println("This is the id from jwt: ", userId)

	var config models.Config
	if err := configs.DB.Where("id = ?", 1).First(&config).Error; err != nil {
		log.Errorf("Failed to get config: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	if config.Pengumuman == false {
		return jsonResponse(c, http.StatusBadRequest, false, "Pengumuman is closed", nil)
	}

	var user models.Datapokok
	if err := configs.DB.Where("user_id = ?", userId).First(&user).Error; err != nil {
		log.Errorf("Failed to get user with user_id %d: %s", userId, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	var nilai models.Nilai
	if err := configs.DB.Where("id = ?", user.ID).First(&nilai).Error; err != nil {
		log.Errorf("Failed to get nilai with id %d: %s", user.ID, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success",
		constans.DATA:    nilai,
		"redirect_wa":    config.RedirectWA,
	})
}
