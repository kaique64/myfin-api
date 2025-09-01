package main

import (
	"log"
	"myfin-api/internal/config"
	"myfin-api/internal/db"
	handlers "myfin-api/internal/handler"
	"myfin-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db.Connect(cfg)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	r.POST("/cash-handling", func(c *gin.Context) {
		handler := handlers.NewCashHandlingHandler(services.NewCashHandlingService())
		
		entry := handler.Save(c)

		c.JSON(200, entry)
	})

	log.Println("🚀 Servidor rodando em http://localhost:8080")
	r.Run(":8080")
}
