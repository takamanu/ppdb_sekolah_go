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
	// Get the search parameter
	paginationParams := ParsePaginationParams(c)

	// Get the search query parameter from the request
	searchQuery := c.QueryParam("search")

	// Create a query builder
	query := configs.DB.Model(&models.User{})

	// Apply the search condition if a search query is provided
	if searchQuery != "" {
		query = query.Where("name LIKE ?", "%"+searchQuery+"%")
	}

	// Apply the pagination parameters
	query = query.Limit(paginationParams.Limit).Offset(paginationParams.Limit * (paginationParams.Page - 1))

	// Preload the "nilai" association
	// Preload the "datapokok" and "nilai" associations
	query = query.Preload("Datapokok.Nilai")

	// Get the paginated results
	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		log.Errorf("Failed to get datapokok: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	for i := range users {
		users[i].Password = "******"
	}

	// Return the paginated users with proper datapokok and nilai associations
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success get all users",
		constans.DATA:    users,
	})
}

func GetUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Invalid id",
		})
	}

	var user models.User
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get user with id %d: %s", id, err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: err.Error(),
		})
	}

	var datapokok []models.Datapokok
	if err := configs.DB.Where("user_id = ?", user.ID).Find(&datapokok).Error; err != nil {
		log.Errorf("Failed to get datapokok for user with id %d: %s", user.ID, err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: err.Error(),
		})
	}

	// Fetch associated nilai data for each datapokok
	for i := range datapokok {
		var nilai []models.Nilai
		if err := configs.DB.Where("datapokok_id = ?", datapokok[i].ID).Find(&nilai).Error; err != nil {
			log.Errorf("Failed to get nilai for datapokok with id %d: %s", datapokok[i].ID, err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				constans.SUCCESS: false,
				constans.MESSAGE: err.Error(),
			})
		}
		datapokok[i].Nilai = nilai
	}

	// Set the user's datapokok
	user.Datapokok = datapokok

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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: err.Error(),
		})
	}

	// Validate all required fields
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "All fields except role must be filled",
		})
	}

	// Validate the email address
	if !validateEmail(user.Email) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Invalid email",
		})
	}

	if len(user.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Password must be at least 8 characters long",
		})
	}

	if IsEmailRegistered(user.Email) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Email address is already registered",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Failed to hash password",
		})
	}

	user.Password = string(hashedPassword)
	loger.Println(user)

	if err := configs.DB.Create(&user).Error; err != nil {
		log.Errorf("Failed to create user: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: err.Error(),
		})
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "Invalid id",
		})
	}

	var user models.User
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get user with id %d: %v", id, err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			constans.SUCCESS: false,
			constans.MESSAGE: "User not found",
		})
	}

	// Check if the datapokok exists
	var datapokok models.Datapokok
	if err := configs.DB.Where("user_id = ?", user.ID).First(&datapokok).Error; err != nil {
		// If the datapokok doesn't exist, just delete the user
		if err == gorm.ErrRecordNotFound {
			if err := configs.DB.Delete(&user).Error; err != nil {
				log.Errorf("Failed to delete user with id %d: %v", id, err)
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					constans.SUCCESS: false,
					constans.MESSAGE: "Failed to delete user id",
				})
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
				return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete datapokok user", nil)
			}

			if err := configs.DB.Delete(&user).Error; err != nil {
				log.Errorf("Failed to delete user with id %d: %v", id, err)
				return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete user", nil)
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
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete nilai user", nil)

	}

	if err := configs.DB.Delete(&datapokok).Error; err != nil {
		log.Errorf("Failed to delete datapokok user with id %d: %v", id, err)
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete datapokok user", nil)
	}

	if err := configs.DB.Delete(&user).Error; err != nil {
		log.Errorf("Failed to delete user with id %d: %v", id, err)
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete user", nil)
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
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid user id", nil)
	}

	// get user by id
	var user models.User
	if err := configs.DB.First(&user, userId).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, "User not found", nil)
	}

	// Validate email
	email := c.FormValue("email")
	if email != "" && !validateEmail(email) {
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid email", nil)
	}

	// Validate password
	newPassword := c.FormValue("password")
	if newPassword != "" && len(newPassword) < 8 {
		return jsonResponse(c, http.StatusBadRequest, false, "Password must be at least 8 characters long", nil)
	}

	if IsEmailRegistered(email) {
		return jsonResponse(c, http.StatusBadRequest, false, "Email address is already registered", nil)

	}

	// bind request body to user struct
	if err := c.Bind(&user); err != nil {
		return jsonResponse(c, http.StatusInternalServerError, false, err.Error(), nil)

	}

	// update password if new password is provided
	if newPassword != "" {
		// encrypt new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			// return jsonResponse(c, http.StatusInternalServerError, false, "Failed to encrypt password")
			return jsonResponse(c, http.StatusInternalServerError, false, "Failed to encrypt password", nil)

		}
		user.Password = string(hashedPassword)
	}

	// save user to database
	if err := configs.DB.Save(&user).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
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

	return jsonResponse(c, http.StatusOK, true, "User successful login", userResponse)
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

func jsonResponse(c echo.Context, status int, success bool, message string, data interface{}) error {
	return c.JSON(status, map[string]interface{}{
		constans.SUCCESS: success,
		constans.MESSAGE: message,
		constans.DATA:    data,
	})
}
