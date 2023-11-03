package configs

import (
	"os"
	"ppdb_sekolah_go/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitGCB(t *testing.T) {
	client, bucketName, err := InitGCB()

	assert.Nil(t, err, "InitGCB() returned an error")
	assert.NotNil(t, client, "InitGCB() returned a nil storage client")
	assert.NotEmpty(t, bucketName, "InitGCB() returned an empty bucket name")

	if client != nil {
		err := client.Close()
		assert.NoError(t, err, "Failed to close the storage client")
	}
}

func TestInitDB(t *testing.T) {
	InitDB()

	assert.NotNil(t, DB, "InitDB() did not initialize the database connection")

	var existingConfig models.Config
	result := DB.First(&existingConfig, 1)
	if result.Error == gorm.ErrRecordNotFound {
	} else if result.Error != nil {
		t.Errorf("Error fetching existing record: %v", result.Error)
	}

	sqlDB, err := DB.DB()
	assert.NoError(t, err, "Failed to get the underlying database connection")
	assert.NoError(t, sqlDB.Close(), "Failed to close the database connection")
}

func TestMain(m *testing.M) {

	exitCode := m.Run()

	os.Exit(exitCode)
}
