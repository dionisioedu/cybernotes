package main

import (
	"log"

	"github.com/dionisioedu/cybernotes/backend/controllers"
	"github.com/dionisioedu/cybernotes/backend/database"
	"github.com/dionisioedu/cybernotes/backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	database.Connect()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/notes", controllers.CreateNote)
	protected.GET("/notes", controllers.GetNotes)
	protected.PUT("/notes/:id", controllers.UpdateNote)
	protected.DELETE("/notes/:id", controllers.DeleteNote)

	r.Run(":8080")
}
