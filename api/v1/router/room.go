package router

import (
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func RoomRouter(router *gin.Engine, roomC controller.RoomController, jwtS service.JWTService) {
	roomRoutes := router.Group("/api/v1/rooms")
	{
		roomRoutes.GET("", middleware.Authenticate(jwtS, constant.EnumRoleUser), roomC.GetAllUserRooms)
		roomRoutes.GET("/:room_id", middleware.Authenticate(jwtS, constant.EnumRoleUser), roomC.GetRoomAndChatsByID)
		roomRoutes.POST("", middleware.Authenticate(jwtS, constant.EnumRoleUser), roomC.CreateNewRoom)
		roomRoutes.DELETE("/:room_id", middleware.Authenticate(jwtS, constant.EnumRoleUser), roomC.DeleteRoomByID)
	}
}
