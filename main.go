package main

import (
	"github.com/gin-gonic/gin"
	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"startup/handler"
	"startup/user"
	"gorm.io/driver/sqlite"
)

func main() {
	// dsn := "mysql:mysql@tcp(127.0.0.1:3306)/startup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	// userService.SaveAvatar(3,"Test-avatar-save.png")

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()

}
