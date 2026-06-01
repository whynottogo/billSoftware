package main

import (
	"log"
	"os"
	_ "time/tzdata"

	"billsoftware/backend/internal/config"
	"billsoftware/backend/internal/database"
	"billsoftware/backend/internal/router"
	"billsoftware/backend/internal/storage"
)

func main() {
	configPath := os.Getenv("BILL_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/app.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	engine, err := database.NewMySQLEngine(cfg.Database)
	if err != nil {
		log.Fatalf("init database failed: %v", err)
	}
	defer func() {
		if err := engine.Close(); err != nil {
			log.Printf("close database failed: %v", err)
		}
	}()

	objectStorage, err := storage.NewObjectStorage(cfg.MinIO)
	if err != nil {
		log.Fatalf("init object storage failed: %v", err)
	}

	httpRouter := router.NewHTTPRouter(cfg, engine, objectStorage)

	if err := httpRouter.Run(cfg.Server.Address()); err != nil {
		log.Fatalf("start http server failed: %v", err)
	}
}
