package constans

import (
	"os"
)

const (
	SECRET_JWT        = "123"
	SUCCESS    string = "success"
	DATA       string = "data"
	MESSAGE    string = "message"
	ERROR      string = "error"
)

var API_KEY string
var DB_USERNAME string
var DB_PORT string
var DB_PASSWORD string
var DB_DATABASE string
var DB_HOST string

func init() {

	API_KEY = os.Getenv("API_KEY")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PORT = os.Getenv("DB_PORT")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_DATABASE = os.Getenv("DB_DATABASE")
	DB_HOST = os.Getenv("DB_HOST")
}
