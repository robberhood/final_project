package main

import (
	"log"

	"github.com/robberhood/final_project/config"
	"github.com/robberhood/final_project/pkg/db"
	"github.com/robberhood/final_project/pkg/server"
)

func main() {
	cfg, err := config.Init("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	dbConnection, err := db.Init(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	log.Println("server started in: ", cfg.Port)
	server.Start(cfg.Port)

}
