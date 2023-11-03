package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"ppdb_sekolah_go/constans"
	"ppdb_sekolah_go/models"
	"time"

	"cloud.google.com/go/storage"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DataSources struct {
	StorageClient *storage.Client
}

func InitGCB() (*storage.Client, string, error) {
	ctx := context.Background()

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "configs/deploy-api-phyton.json")
	// Sets your Google Cloud Platform project ID.
	projectID := "deploy-api-phyton"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, "", err
	}

	// Sets the name for the bucket.
	bucketName := "data-sekolah-app"

	// Check if the bucket already exists.
	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		// If the bucket doesn't exist, create it.
		// Note: You can handle the error here if needed.
		log.Printf("Bucket %v does not exist, creating it...\n", bucketName)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		if err := bucket.Create(ctx, projectID, nil); err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
			return nil, "", err
		}
		fmt.Printf("Bucket %v created.\n", bucketName)
	} else {
		fmt.Printf("Usinga this test existing lolo bucket: %v\n", bucketName)
	}

	return client, bucketName, nil
}

func InitDB() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", constans.DB_USERNAME, constans.DB_PASSWORD, constans.DB_HOST, constans.DB_PORT, constans.DB_DATABASE)
	// dsn := fmt.Sprintf("root:@tcp(127.0.0.1:3306)/ppdb_smp?charset=utf8mb4&parseTime=True&loc=Local")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Nilai{})
	DB.AutoMigrate(&models.Datapokok{})
	DB.AutoMigrate(&models.Config{})

	var existingConfig models.Config
	if result := DB.First(&existingConfig, 1); result.Error == gorm.ErrRecordNotFound {
		// Record with ID = 1 doesn't exist, so create it
		config := models.Config{
			ID:         1,
			Pengumuman: true,
			RedirectWA: "https://example.com/whatsapp",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Insert the sample data into the database
		result := DB.Create(&config)
		if result.Error != nil {
			fmt.Println(result.Error)
			return
		}

		fmt.Println("Seeder executed successfully.")
	} else if result.Error != nil {
		// Handle the error if there was a problem fetching the existing record
		fmt.Println(result.Error)
	} else {
		fmt.Println("Record with ID = 1 already exists, no need to seed.")
	}

}
