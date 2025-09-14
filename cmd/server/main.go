package main

import (
	"fmt"
	"log"

	"myfin-api/internal/config"
	"myfin-api/internal/db"
	handlers "myfin-api/internal/handler"
	"myfin-api/internal/repository"
	"myfin-api/internal/services"

	"github.com/gin-gonic/gin"
)

const transactionsPath = "/transactions"

var transactionsIDPath = fmt.Sprintf("%s/:id", transactionsPath)

func main() {
	cfg := config.LoadConfig()

	db.Connect(cfg)

	r := gin.Default()
	handler := handlers.NewTransactionsHandler(services.NewTransactionsService(repository.NewTransactionsEntryRepository(db.MongoDatabase)))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.POST(transactionsPath, func(c *gin.Context) {
		handler.Save(c)
	})

	r.GET(transactionsPath, func(c *gin.Context) {
		handler.GetAll(c)
	})

	r.GET(transactionsIDPath, func(c *gin.Context) {
		handler.GetByID(c)
	})

	r.PUT(transactionsIDPath, func(c *gin.Context) {
		handler.Update(c)
	})

	r.DELETE(transactionsIDPath, func(c *gin.Context) {
		handler.Delete(c)
	})

	log.Println("🚀 Servidor rodando em http://localhost:8080")
	r.Run(":8080")
}