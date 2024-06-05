package main

import (
	"log"

	cloud "github.com/JIIL07/cloudFiles-manager/client"
	server "github.com/JIIL07/cloudFiles-manager/server"
)

var sqlite *cloud.SQLiteDB

func main() {
	db, err := sqlite.PrepareLocalDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
