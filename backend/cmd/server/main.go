package main

import (
	"log"

	"billsoftware/backend/internal/config"
	"billsoftware/backend/internal/database"
	"billsoftware/backend/internal/router"
)

func main() {
	cfg, err := config.LoadConfig("configs/app.yaml")
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

	httpRouter := router.NewHTTPRouter(cfg, engine)

	if err := httpRouter.Run(cfg.Server.Address()); err != nil {
		log.Fatalf("start http server failed: %v", err)
	}
}

