package main

import (
	"backend-gin-gonic/config"
	"backend-gin-gonic/database"
	"backend-gin-gonic/handlers"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

// @title           Swagger Auth API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func GetRoute(r *gin.Engine) {

	r.POST("/api/login", handlers.Login)
	r.GET("/api/refresh", handlers.Refresh)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func init() {
	config.LoadEnvs()
	database.ConnectToDB()
}

func main() {

	//database.ConnectToDB()

	fmt.Println("Service Running")

	if database.DB == nil {
		log.Fatal("Database connection is not initialized")
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Update with your Swagger UI origin
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	GetRoute(router)

	router.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}
