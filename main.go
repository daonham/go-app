package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/daonham/go-app/controllers"
)

func main() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(gzip.Gzip(gzip.BestCompression))

	router.GET("/post", controllers.GetPosts)
	router.GET("/post/:id", controllers.GetPost)
	router.POST("/post", controllers.CreatePost)
	router.PUT("/post/:id", controllers.UpdatePost)
	router.DELETE("/post/:id", controllers.DeletePost)

	router.Run("localhost:8070")
}
