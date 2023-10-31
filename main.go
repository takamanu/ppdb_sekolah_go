package main

import (
	"os"
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/routes"
)

func main() {
	configs.InitDB()
	configs.InitGCB()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000" // Default port if not specified in the environment variable
	}

	e := routes.New()

	e.Logger.Fatal(e.Start(":" + port))
}
