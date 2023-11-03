package controllers

import (
	"context"
	"errors"
	"fmt"
	"io"
	loger "log"
	"net/http"
	"os"
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/constans"
	"ppdb_sekolah_go/models"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// func GetDatapokokController(c echo.Context) error {

// 	var users []models.Datapokok
// 	if err := configs.DB.Find(&users).Error; err != nil {
// 		log.Errorf("Failed to get datapokok: %s", err.Error())
// 		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

// 	}
// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		constans.SUCCESS: true,
// 		constans.MESSAGE: "Success get all datapokok",
// 		constans.DATA:    users,
// 	})
// }

// func SearchDatapokokController(c echo.Context) error {
// 	// Parse pagination parameters
// 	paginationParams := ParsePaginationParams(c)

// 	// Get the search query parameter for name
// 	nameQuery := c.QueryParam("name")

// 	// Initialize the query with the search and pagination parameters
// 	query := configs.DB.Model(&models.Datapokok{})

// 	// Apply the name search filter if the nameQuery is provided
// 	if nameQuery != "" {
// 		query = query.Where("nama_lengkap LIKE ?", "%"+nameQuery+"%")
// 	}

// 	var users []models.Datapokok

// 	// Query the database with the search and pagination parameters
// 	result, err := GetPaginatedData(c, query, paginationParams, users)
// 	if err != nil {
// 		log.Errorf("Failed to get datapokok: %s", err.Error())
// 		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		constans.SUCCESS: true,
// 		constans.MESSAGE: "Success get datapokok",
// 		constans.DATA:    result,
// 	})
// }

func GetDatapokokController(c echo.Context) error {
	// Parse pagination parameters
	paginationParams := ParsePaginationParams(c)

	// Get the search query parameter from the request
	searchQuery := c.QueryParam("search")

	// Create a query builder
	query := configs.DB.Model(&models.Datapokok{})

	// Apply the search condition if a search query is provided
	if searchQuery != "" {
		query = query.Where("nama_lengkap LIKE ?", "%"+searchQuery+"%")
	}

	// Apply the pagination parameters
	query = query.Limit(paginationParams.Limit).Offset(paginationParams.Limit * (paginationParams.Page - 1))

	// Preload the "nilai" association
	query = query.Preload("Nilai")

	// Get the paginated results
	var datapokokList []models.Datapokok
	if err := query.Find(&datapokokList).Error; err != nil {
		log.Errorf("Failed to get datapokok: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)
	}

	// Return the results
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success get datapokok",
		constans.DATA:    datapokokList,
	})
}

// Use the same approach for other Get API functions

func GetDatapokokControllerByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid id", nil)

	}

	var user models.Datapokok
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get datapokok with id %d: %s", id, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	var nilai models.Nilai
	if err := configs.DB.Where("datapokok_id = ?", id).First(&nilai).Error; err != nil {
		log.Errorf("Failed to get nilai with datapokok_id %d: %s", id, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	user.Nilai = append(user.Nilai, nilai)

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success get datapokok by ID",
		constans.DATA:    user,
	})
}

func CreateDatapokokController(c echo.Context, client *storage.Client, bucketName string) error {
	// Create a request structure that includes Datapokok and Nilai data

	requestData := struct {
		Datapokok models.Datapokok `json:"datapokok"`
		Nilai     models.Nilai     `json:"nilai"`
	}{}

	// Bind the request data from the JSON body
	if err := c.Bind(&requestData); err != nil {
		log.Errorf("Failed to bind request: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	userIDDatapokokStr := c.FormValue("user_id")
	userIDDatapokok, err := strconv.ParseUint(userIDDatapokokStr, 10, 0)
	if err != nil {
		log.Errorf("Failed to convert user_id to a uint: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid user_id", nil)
	}

	requestData.Datapokok.UserID = uint64(userIDDatapokok)

	requestData.Datapokok.Email = c.FormValue("email")
	requestData.Datapokok.NamaLengkap = c.FormValue("nama_lengkap")
	requestData.Datapokok.NISN = c.FormValue("nisn")
	requestData.Datapokok.JenisKelamin = c.FormValue("jenis_kelamin")
	requestData.Datapokok.TempatLahir = c.FormValue("tempat_lahir")

	if IsEmailRegisteredDatapokok(requestData.Datapokok.Email) {
		return jsonResponse(c, http.StatusBadRequest, false, "Email address is already registered", nil)
	}

	if IsNISNRegisteredDatapokok(requestData.Datapokok.NISN) {
		return jsonResponse(c, http.StatusBadRequest, false, "NISN is already registered", nil)
	}

	// Date of birth handling
	dobStr := c.FormValue("tanggal_lahir")
	dob, err := time.Parse("2006-01-02", dobStr)
	if err == nil {
		requestData.Datapokok.TanggalLahir = &dob
	}

	requestData.Datapokok.AsalSekolah = c.FormValue("asal_sekolah")
	requestData.Datapokok.NamaAyah = c.FormValue("nama_ayah")
	requestData.Datapokok.NoWaAyah = c.FormValue("no_wa_ayah")
	requestData.Datapokok.NamaIbu = c.FormValue("nama_ibu")
	requestData.Datapokok.NoWaIbu = c.FormValue("no_wa_ibu")
	requestData.Datapokok.Jurusan = c.FormValue("jurusan")

	if err := ValidateDatapokokFields(requestData.Datapokok); err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// Handle file upload
	image, err := c.FormFile("pas_foto")
	if err != nil {
		log.Errorf("Failed to get the image file: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, "Image upload failed", nil)
	}

	// Generate a unique filename using a UUID
	uniqueFilename := uuid.NewString()

	// Upload the image to the existing Google Cloud Storage bucket
	ctx := context.Background()
	wc := client.Bucket(bucketName).Object(uniqueFilename).NewWriter(ctx)
	defer wc.Close()

	src, err := image.Open()
	if err != nil {
		log.Errorf("Failed to open the image file: %s", err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to process image", nil)

	}
	defer src.Close()

	if _, err = io.Copy(wc, src); err != nil {
		log.Errorf("Failed to copy the image to the bucket: %s", err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to upload image", nil)
	}

	requestData.Datapokok.PasFoto = "https://storage.googleapis.com/" + bucketName + "/" + uniqueFilename

	// Create the Datapokok record in the database
	if err := configs.DB.Create(&requestData.Datapokok).Error; err != nil {
		log.Errorf("Failed to create datapokok: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// Now requestData.Datapokok.ID contains the ID of the newly created Datapokok record
	loger.Println("Created Datapokok with ID:", requestData.Datapokok.ID)

	// Set the Nilai's DatapokokID to the ID of the created Datapokok record
	requestData.Nilai.DataPokokID = requestData.Datapokok.ID
	requestData.Nilai.BahasaIndonesia = 0
	requestData.Nilai.IlmuPengetahuanAlam = 0
	requestData.Nilai.Matematika = 0
	requestData.Nilai.TestMembacaAlQuran = 0
	requestData.Nilai.Status = "BELUM LULUS"

	// Create the Nilai record in the database
	if err := configs.DB.Create(&requestData.Nilai).Error; err != nil {
		log.Errorf("Failed to create nilai: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// requestData.Nilai.Utama

	requestData.Datapokok.Nilai = append(requestData.Datapokok.Nilai, requestData.Nilai)

	// Return a response
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success create new Datapokok and Nilai",
		constans.DATA:    requestData.Datapokok,
	})
}

// delete user by id
func DeleteDatapokokController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("Invalid id: %s", c.Param("id"))
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid id", nil)
	}

	var user models.Datapokok
	if err := configs.DB.First(&user, id).Error; err != nil {
		log.Errorf("Failed to get datapokok with id %d: %v", id, err)
		return jsonResponse(c, http.StatusNotFound, false, "User not found", nil)

	}

	if err := configs.DB.Delete(&user).Error; err != nil {
		log.Errorf("Failed to delete datapokok with id %d: %v", id, err)
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to delete datapokok", nil)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "success deleted datapokok",
	})
}

// update user by id
func UpdateDatapokokController(c echo.Context, client *storage.Client, bucketName string) error {
	// get user id from url param
	// get user id from url param
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid datapokok id", nil)
	}

	// get user by id
	var user models.Datapokok
	if err := configs.DB.First(&user, userId).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, "Datapokok not found", nil)
	}

	if c.FormValue("user_id") != "" {
		userIDDatapokokStr := c.FormValue("user_id")
		userIDDatapokok, err := strconv.ParseUint(userIDDatapokokStr, 10, 0)
		if err != nil {
			log.Errorf("Failed to convert user_id to a uint: %s", err.Error())
			return jsonResponse(c, http.StatusBadRequest, false, "Invalid user_id", nil)
		}

		user.UserID = uint64(userIDDatapokok)
	}

	if c.FormValue("email") != "" {
		user.Email = c.FormValue("email")
		if !isEmailValid(user.Email) {
			return errors.New("invalid email address")
		}
	}

	if c.FormValue("nama_lengkap") != "" {
		user.NamaLengkap = c.FormValue("nama_lengkap")
	}

	if c.FormValue("nisn") != "" {
		user.NISN = c.FormValue("nisn")
	}

	if c.FormValue("jenis_kelamin") != "" {
		user.JenisKelamin = c.FormValue("jenis_kelamin")
	}

	if c.FormValue("tempat_lahir") != "" {
		user.TempatLahir = c.FormValue("tempat_lahir")
	}

	if c.FormValue("tanggal_lahir") != "" {
		dobStr := c.FormValue("tanggal_lahir")
		dob, err := time.Parse("2006-01-02", dobStr)
		if err == nil {
			user.TanggalLahir = &dob
		}
	}
	// Date of birth handling

	if c.FormValue("asal_sekolah") != "" {
		user.AsalSekolah = c.FormValue("asal_sekolah")
	}

	if c.FormValue("nama_ayah") != "" {
		user.NamaAyah = c.FormValue("nama_ayah")
	}

	if c.FormValue("no_wa_ayah") != "" {
		user.NoWaAyah = c.FormValue("no_wa_ayah")
	}

	if c.FormValue("nama_ibu") != "" {
		user.NamaIbu = c.FormValue("nama_ibu")
	}

	if c.FormValue("no_wa_ibu") != "" {
		user.NoWaIbu = c.FormValue("no_wa_ibu")
	}

	if c.FormValue("jurusan") != "" {
		user.Jurusan = c.FormValue("jurusan")
	}

	// // Create the Datapokok record in the database
	// if err := configs.DB.Create(&requestData.Datapokok).Error; err != nil {
	// 	log.Errorf("Failed to create datapokok: %s", err.Error())
	// 			return jsonResponse(c, http.StatusBadRequest, false,  err.Error(), nil)

	// }

	// Handle file upload

	image, err := c.FormFile("pas_foto")
	if err == nil {
		// Generate a unique filename using a UUID
		uniqueFilename := uuid.NewString()

		// Upload the image to the existing Google Cloud Storage bucket
		ctx := context.Background()
		wc := client.Bucket(bucketName).Object(uniqueFilename).NewWriter(ctx)
		defer wc.Close()

		src, err := image.Open()
		if err != nil {
			log.Errorf("Failed to open the image file: %s", err.Error())
			return jsonResponse(c, http.StatusInternalServerError, false, "Failed to process image", nil)
		}
		defer src.Close()

		if _, err = io.Copy(wc, src); err != nil {
			log.Errorf("Failed to copy the image to the bucket: %s", err.Error())
			return jsonResponse(c, http.StatusInternalServerError, false, "Failed to upload image", nil)
		}

		user.PasFoto = "https://storage.googleapis.com/" + bucketName + "/" + uniqueFilename
	}

	if c.Request().Method == "PUT" && user.Nilai != nil {
		return jsonResponse(c, http.StatusForbidden, false, "You cant update user nilai", nil)

	}

	// // validate user fields
	// if err := ValidateDatapokokFields(user); err != nil {
	// 			return jsonResponse(c, http.StatusBadRequest, false,  err.Error(), nil)

	// }

	// update user to database
	if err := configs.DB.Save(&user).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	var nilai models.Nilai
	if err := configs.DB.Where("datapokok_id = ?", user.ID).First(&nilai).Error; err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, "Datapokok not found", nil)
	}

	user.Nilai = append(user.Nilai, nilai)

	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success datapokok updated",
		constans.DATA:    user,
	})
}

func GetDatapokokControllerSiswa(c echo.Context) error {

	userId := c.Get("userId")
	fmt.Println("This is the id from jwt: ", userId)

	// GET THE USER FROM THE DATABASE

	var user models.Datapokok
	if err := configs.DB.Where("user_id = ?", userId).First(&user).Error; err != nil {
		log.Errorf("Failed to get user with user_id %d: %s", userId, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// GET THE NILAIS FROM THE DATABASE
	var nilai models.Nilai
	if err := configs.DB.Where("datapokok_id = ?", user.ID).First(&nilai).Error; err != nil {
		log.Errorf("Failed to get nilai with datapokok_id %d: %s", userId, err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// APPEND THE NILAIS TO THE USER
	user.Nilai = append(user.Nilai, nilai)

	// RETURN THE USER AS JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success get datapokok by ID",
		constans.DATA:    user,
	})
}

func CreateDatapokokControllerSiswa(c echo.Context, client *storage.Client, bucketName string) error {

	userId, ok := c.Get("userId").(float64)
	if !ok {
		return jsonResponse(c, http.StatusBadRequest, false, "Invalid user ID", nil)
	}

	// Create a request structure that includes Datapokok and Nilai data
	requestData := struct {
		Datapokok models.Datapokok `json:"datapokok"`
		Nilai     models.Nilai     `json:"nilai"`
	}{}

	// Bind the request data from the JSON body
	if err := c.Bind(&requestData); err != nil {
		log.Errorf("Failed to bind request: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	requestData.Datapokok.UserID = uint64(userId)

	requestData.Datapokok.Email = c.FormValue("email")
	requestData.Datapokok.NamaLengkap = c.FormValue("nama_lengkap")
	requestData.Datapokok.NISN = c.FormValue("nisn")
	requestData.Datapokok.JenisKelamin = c.FormValue("jenis_kelamin")
	requestData.Datapokok.TempatLahir = c.FormValue("tempat_lahir")

	if IsEmailRegisteredDatapokok(requestData.Datapokok.Email) {
		return jsonResponse(c, http.StatusBadRequest, false, "Email address is already registered", nil)
	}

	if IsNISNRegisteredDatapokok(requestData.Datapokok.NISN) {
		return jsonResponse(c, http.StatusBadRequest, false, "NISN is already registered", nil)
	}

	// Date of birth handling
	dobStr := c.FormValue("tanggal_lahir")
	dob, err := time.Parse("2006-01-02", dobStr)
	if err == nil {
		requestData.Datapokok.TanggalLahir = &dob
	}

	requestData.Datapokok.AsalSekolah = c.FormValue("asal_sekolah")
	requestData.Datapokok.NamaAyah = c.FormValue("nama_ayah")
	requestData.Datapokok.NoWaAyah = c.FormValue("no_wa_ayah")
	requestData.Datapokok.NamaIbu = c.FormValue("nama_ibu")
	requestData.Datapokok.NoWaIbu = c.FormValue("no_wa_ibu")
	requestData.Datapokok.Jurusan = c.FormValue("jurusan")

	if err := ValidateDatapokokFields(requestData.Datapokok); err != nil {
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// Create the Datapokok record in the database
	if err := configs.DB.Create(&requestData.Datapokok).Error; err != nil {
		log.Errorf("Failed to create datapokok: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// Handle file upload
	image, err := c.FormFile("pas_foto")
	if err != nil {
		log.Errorf("Failed to get the image file: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, "Image upload failed", nil)
	}

	// Generate a unique filename using a UUID
	uniqueFilename := uuid.NewString()

	// Upload the image to the existing Google Cloud Storage bucket
	ctx := context.Background()
	wc := client.Bucket(bucketName).Object(uniqueFilename).NewWriter(ctx)
	defer wc.Close()

	src, err := image.Open()
	if err != nil {
		log.Errorf("Failed to open the image file: %s", err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to process image", nil)
	}
	defer src.Close()

	if _, err = io.Copy(wc, src); err != nil {
		log.Errorf("Failed to copy the image to the bucket: %s", err.Error())
		return jsonResponse(c, http.StatusInternalServerError, false, "Failed to upload image", nil)
	}

	requestData.Datapokok.PasFoto = "https://storage.googleapis.com/" + bucketName + "/" + uniqueFilename

	// Now requestData.Datapokok.ID contains the ID of the newly created Datapokok record
	loger.Println("Created Datapokok with ID:", requestData.Datapokok.ID)

	// Set the Nilai's DatapokokID to the ID of the created Datapokok record
	requestData.Nilai.DataPokokID = requestData.Datapokok.ID
	requestData.Nilai.BahasaIndonesia = 0
	requestData.Nilai.IlmuPengetahuanAlam = 0
	requestData.Nilai.Matematika = 0
	requestData.Nilai.TestMembacaAlQuran = 0
	requestData.Nilai.Status = "BELUM LULUS"

	// Create the Nilai record in the database
	if err := configs.DB.Create(&requestData.Nilai).Error; err != nil {
		log.Errorf("Failed to create nilai: %s", err.Error())
		return jsonResponse(c, http.StatusBadRequest, false, err.Error(), nil)

	}

	// requestData.Nilai.Utama

	requestData.Datapokok.Nilai = append(requestData.Datapokok.Nilai, requestData.Nilai)

	// Return a response
	return c.JSON(http.StatusOK, map[string]interface{}{
		constans.SUCCESS: true,
		constans.MESSAGE: "Success create new Datapokok and Nilai",
		constans.DATA:    requestData.Datapokok,
	})
}

func IsEmailRegisteredDatapokok(email string) bool {
	var user models.Datapokok
	if err := configs.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}

func IsNISNRegisteredDatapokok(nisn string) bool {
	var user models.Datapokok
	if err := configs.DB.Where("nisn = ?", nisn).First(&user).Error; err != nil {
		return false
	}
	return true
}

func ValidateDatapokokFields(datapokok models.Datapokok) error {
	// Validate email
	if !isEmailValid(datapokok.Email) {
		return errors.New("invalid email address")
	}
	// Validate pasfoto
	// if !isPasfotoValid(datapokok.PasFoto) {
	// 	return errors.New("pasfoto must be an image file")
	// }

	if datapokok.AsalSekolah == "" {
		return errors.New("Asal sekolah is required")
	}

	if datapokok.JenisKelamin == "" {
		return errors.New("Gender is required")
	}
	if datapokok.NISN == "" {
		return errors.New("NISN is required")
	}
	if datapokok.NamaAyah == "" {
		return errors.New("Nama ayah is required")
	}
	if datapokok.NamaIbu == "" {
		return errors.New("Nama ibu is required")
	}
	if datapokok.NamaLengkap == "" {
		return errors.New("Nama Lengkap is required")
	}
	if datapokok.NoWaAyah == "" {
		return errors.New("No WA Ayah is required")
	}
	if datapokok.NoWaIbu == "" {
		return errors.New("No WA Ibu is required")
	}

	if datapokok.TempatLahir == "" {
		return errors.New("Tempat Lahir is required")
	}

	if datapokok.TanggalLahir == nil {
		return errors.New("Tanggal lahir is required")
	}

	// Validate jurusan
	if datapokok.Jurusan == "" {
		return errors.New("jurusan is required")
	}

	return nil
}

func isEmailValid(email string) bool {
	// Use a regular expression to validate the email address
	re := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return re.MatchString(email)
}

func isPasfotoValid(pasfoto string) bool {
	// Check if the pasfoto is an image file
	file, err := os.Open(pasfoto)
	if err != nil {
		return false
	}
	defer file.Close()

	// Determine the MIME type of the file
	buffer := make([]byte, 512) // Read the first 512 bytes to detect the MIME type
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	mimetype := http.DetectContentType(buffer)

	// A valid pasfoto must be an image file of type png, jpeg, or gif
	return mimetype == "image/png" || mimetype == "image/jpg" || mimetype == "image/jpeg" || mimetype == "image/gif"
}
