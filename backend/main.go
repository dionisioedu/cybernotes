package main

import (
	"github.com/dionisioedu/cybernotes/backend/controllers"
	"github.com/dionisioedu/cybernotes/backend/database"
	"github.com/dionisioedu/cybernotes/backend/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
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
	protected.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(200, gin.H{"message": "Acesso permitido", "user_id": userID})
	})

	r.Run(":8080")
}
