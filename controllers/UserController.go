package controllers

import (
	"fmt"
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/constans"
	m "ppdb_sekolah_go/middlewares"
	"ppdb_sekolah_go/models"
	"regexp"

	loger "log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsersController(c echo.Context) error {
	var users []models.User
	if err := configs.DB.Find(&users).Error; err != nil {
		log.Errorf("Failed to get users: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success get all users",
		constans.DATA:    users,
		//USAGE OF THE GLOBAL VARIABLE
	})
}

func GetUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	var user models.User
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get user with id %d: %s", id, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success get user by id",
		constans.DATA:    user,
	})
}

func CreateUserController(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Failed to bind request: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate all required fields
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "All fields except role must be filled")
	}

	// Validate the email address
	if !validateEmail(user.Email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	// Validate the password
	if len(user.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	if IsEmailRegistered(user.Email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	user.Password = string(hashedPassword)
	loger.Println(user)

	if err := configs.DB.Create(&user).Error; err != nil {
		log.Errorf("Failed to create user: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success create new user",
		constans.DATA:    user,
	})
}

// delete user by id
func DeleteUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	var user models.User
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get user with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Check if the datapokok exists
	var datapokok models.Datapokok
	if err := configs.DB.Where("user_id = ?", user.ID).First(&datapokok).Error; err != nil {
		// If the datapokok doesn't exist, just delete the user
		if err == gorm.ErrRecordNotFound {
			if err := configs.DB.Delete(&user).Error; err != nil {
				log.Errorf("Failed to delete user with id %d: %v", id, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				constans.SUCCESS: true,
				constans.MESSAGE: "Success deleted user",
			})
		}

		// If there's another error, return it
		return err
	}

	// Check if the nilai exists
	var nilai models.Nilai
	if err := configs.DB.Where("datapokok_id = ?", datapokok.ID).First(&nilai).Error; err != nil {
		// If the nilai doesn't exist, delete the datapokok and user
		if err == gorm.ErrRecordNotFound {
			if err := configs.DB.Delete(&datapokok).Error; err != nil {
				log.Errorf("Failed to delete datapokok user with id %d: %v", id, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete datapokok user")
			}

			if err := configs.DB.Delete(&user).Error; err != nil {
				log.Errorf("Failed to delete user with id %d: %v", id, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				constans.SUCCESS: true,
				constans.MESSAGE: "Success deleted user and datapokok",
			})
		}

		// If there's another error, return it
		return err
	}

	// If all three exist, delete all three
	if err := configs.DB.Delete(&nilai).Error; err != nil {
		log.Errorf("Failed to delete nilai user with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete nilai user")
	}

	if err := configs.DB.Delete(&datapokok).Error; err != nil {
		log.Errorf("Failed to delete datapokok user with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete datapokok user")
	}

	if err := configs.DB.Delete(&user).Error; err != nil {
		log.Errorf("Failed to delete user with id %d: %v", id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success deleted user, datapokok, and nilai",
	})
}

// update user by id
func UpdateUserController(c echo.Context) error {
	// get user id from url param
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	// get user by id
	var user models.User
	if err := configs.DB.First(&user, userId).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	// Validate email
	email := c.FormValue("email")
	if email != "" && !validateEmail(email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email")
	}

	// Validate password
	newPassword := c.FormValue("password")
	if newPassword != "" && len(newPassword) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	if IsEmailRegistered(email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is already registered")
	}

	// bind request body to user struct
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// update password if new password is provided
	if newPassword != "" {
		// encrypt new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encrypt password")
		}
		user.Password = string(hashedPassword)
	}

	// save user to database
	if err := configs.DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success user updated",
		constans.DATA:    user,
	})
}

func LoginUserController(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	err := configs.DB.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Failed to login",
			constans.ERROR:   err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password"))); err != nil {
		// fmt.Println(err)
		fmt.Println("pass :", c.FormValue("password"))
		fmt.Println("err :", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	}

	fmt.Println("pass :", c.FormValue("password"))

	token, err := m.CreateToken(int(user.ID), user.Name, int(user.Role))
	fmt.Printf("UserID: %v, UserName: %v, UserRole: %v", user.ID, user.Name, user.Role)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Failed to login",
			constans.ERROR:   err.Error(),
		})
	}
	userResponse := models.UserResponse{user.ID, user.Name, user.Email, user.Role, token}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success login",
		constans.DATA:    userResponse,
	})
}

func IsEmailRegistered(email string) bool {
	var user models.User
	if err := configs.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,}$`)
	return re.MatchString(email)
}
