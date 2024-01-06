package main

import (
	"test-app/db"
	"test-app/handlers"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := gin.Default()

	dbConn, err := db.InitDB()
	if err != nil {
		panic(err)
	}

	userHandler := handlers.NewUserHandler(dbConn)

	router.POST("/users", userHandler.CreateUser)
	router.POST("/users/generateotp", userHandler.GenerateOTP)
	router.POST("/users/verifyotp", userHandler.VerifyOTP)

	router.Run(":8080")
}
