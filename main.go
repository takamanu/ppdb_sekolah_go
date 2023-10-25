package main

import (
	"ppdb_sekolah_go/configs"
	"ppdb_sekolah_go/routes"
)

func main() {
	configs.InitDB()

	e := routes.New()

	e.Logger.Fatal(e.Start(":8000"))
}
