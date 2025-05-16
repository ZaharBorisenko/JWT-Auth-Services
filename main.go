package main

import (
	"github.com/ZaharBorisenko/JWT-Auth-Services/database"
	"github.com/ZaharBorisenko/JWT-Auth-Services/models"
	"github.com/ZaharBorisenko/JWT-Auth-Services/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}
	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to database")
	}

	if err := models.Migrate(db); err != nil {
		log.Fatal("Could not migrate models")
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success api-1"})
	})
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success api-2"})
	})

	router.Run(":8080")
}
