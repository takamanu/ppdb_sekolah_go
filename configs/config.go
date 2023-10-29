package configs

import (
	"context"
	"fmt"
	"log"
	"os"
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

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "configs/deploy-api-phyton-c5c9b7d8a6df.json")
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
		fmt.Printf("Using existing bucket: %v\n", bucketName)
	}

	return client, bucketName, nil
}

func InitDB() {

	dsn := "root:P-8VA^=pL2dX`D8=@tcp(35.240.201.186:3306)/ppdb_smp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Nilai{})
	DB.AutoMigrate(&models.Datapokok{})
	DB.AutoMigrate(&models.Config{})

}
