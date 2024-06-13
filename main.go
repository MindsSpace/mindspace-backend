package main

import (
	"fmt"
	"os"

	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/router"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/config"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		db   = config.DBSetup()
		jwtS = service.NewJWTService()

		txR   = repository.NewTxRepository(db)
		userR = repository.NewUserRepository(txR)
		chatR = repository.NewChatRepository(txR)
		roomR = repository.NewRoomRepository(txR)

		userS = service.NewUserService(userR)
		chatS = service.NewChatService(chatR, roomR)
		roomS = service.NewRoomService(roomR, chatR)

		userC = controller.NewUserController(userS, jwtS)
		chatC = controller.NewChatController(chatS)
		roomC = controller.NewRoomController(roomS)
		fileC = controller.NewFileController()
	)

	defer config.DBClose(db)

	// Setting Up Server
	server := gin.Default()
	server.Use(
		middleware.CORSMiddleware(),
	)

	// Setting Up Routes
	router.UserRouter(server, userC, jwtS)
	router.RoomRouter(server, roomC, jwtS)
	router.ChatRouter(server, chatC, jwtS)
	router.FileRouter(server, fileC)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := server.Run(":" + port)
	if err != nil {
		fmt.Println("Server failed to start: ", err)
		return
	}
}
