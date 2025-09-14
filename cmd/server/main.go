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

const cashHandlingPath = "/cash-handling"

var cashHandlingIDPath = fmt.Sprintf("%s/:id", cashHandlingPath)

func main() {
	cfg := config.LoadConfig()

	db.Connect(cfg)

	r := gin.Default()
	handler := handlers.NewCashHandlingHandler(services.NewCashHandlingService(repository.NewCashHandlingEntryRepository(db.MongoDatabase)))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.POST(cashHandlingPath, func(c *gin.Context) {
		handler.Save(c)
	})

	r.GET(cashHandlingPath, func(c *gin.Context) {
		handler.GetAll(c)
	})

	r.GET(cashHandlingIDPath, func(c *gin.Context) {
		handler.GetByID(c)
	})

	r.PUT(cashHandlingIDPath, func(c *gin.Context) {
		handler.Update(c)
	})

	r.DELETE(cashHandlingIDPath, func(c *gin.Context) {
		handler.Delete(c)
	})

	log.Println("🚀 Servidor rodando em http://localhost:8080")
	r.Run(":8080")
}
