package main

import (
	"log"

	"server-blog-v2/config"
	"server-blog-v2/internal/app"
)

// @title Blog API
// @version 2.0
// @description 博客系统 API

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	app.Run(cfg)
}
