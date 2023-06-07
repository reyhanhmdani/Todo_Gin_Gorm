package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"todoGin/config"
	"todoGin/database"
	"todoGin/router"
	"todoGin/service"
)

func setupLogOutput() {
	f, _ := os.Create("gin-log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {

	setupLogOutput()

	ctx := context.Background()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	var cfg config.Config

	// pr
	err1 := envconfig.Process("", &cfg)
	if err1 != nil {
		log.Fatal("error", err1)
	}
	// INITAL DATABASE
	db, err := database.Databaseinit(ctx, &cfg)
	if err != nil {
		return
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Error running schema migration %v", err)
	}

	// initial repo
	todoRepo := database.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	routeBuilder := router.NewRouteBuilder(todoService)
	routeInit := routeBuilder.RouteInit()
	err = routeInit.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
