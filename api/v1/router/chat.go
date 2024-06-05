package router

import (
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func ChatRouter(router *gin.Engine, chatC controller.ChatController, jwtS service.JWTService) {
	chatRoutes := router.Group("/api/v1/chats")
	{
		chatRoutes.GET("/:chat_id", middleware.Authenticate(jwtS, constant.EnumRoleUser), chatC.GetChatByID)
		chatRoutes.POST("", middleware.Authenticate(jwtS, constant.EnumRoleUser), chatC.CreateNewChat)
		chatRoutes.DELETE("/:chat_id", middleware.Authenticate(jwtS, constant.EnumRoleUser), chatC.DeleteChatByID)
	}
}
