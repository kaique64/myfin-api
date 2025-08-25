package main

import (
	"finix-api/internal/config"
	"finix-api/internal/db"
	"log"

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

	log.Println("🚀 Servidor rodando em http://localhost:8080")
	r.Run(":8080")
}
