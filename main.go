package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"

	"github.com/daonham/go-app/controllers"
)

// Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		controllers.TokenValid(c)
		c.Next()
	}
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		panic("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(RequestIDMiddleware())
	router.Use(gzip.Gzip(gzip.BestCompression))

	router.GET("/post", TokenAuthMiddleware(), controllers.GetPosts)
	router.GET("/post/:id", controllers.GetPost)
	router.POST("/post", controllers.CreatePost)
	router.PUT("/post/:id", controllers.UpdatePost)
	router.DELETE("/post/:id", controllers.DeletePost)

	router.GET("/user", controllers.GetUsers)
	router.GET("/user/:id", controllers.GetUser)
	router.POST("/user", controllers.CreateUser)
	router.PUT("/user/:id", controllers.UpdateUser)
	router.DELETE("/user/:id", controllers.DeleteUser)

	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	router.POST("/token/refresh", controllers.RefreshToken)

	router.Run("localhost:8070")
}
