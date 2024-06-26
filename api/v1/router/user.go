package router

import (
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, userC controller.UserController, jwtS service.JWTService) {
	userRoutes := router.Group("/api/v1/users")
	{
		// admin routes
		userRoutes.GET("", middleware.Authenticate(jwtS, constant.EnumRoleAdmin), userC.GetAllUsers)
		userRoutes.PATCH("/:user_id", middleware.Authenticate(jwtS, constant.EnumRoleAdmin), userC.UpdateUserByID)
		userRoutes.DELETE("/:user_id", middleware.Authenticate(jwtS, constant.EnumRoleAdmin), userC.DeleteUserByID)

		// user routes
		userRoutes.GET("/me", middleware.Authenticate(jwtS, constant.EnumRoleUser), userC.GetMe)
		userRoutes.DELETE("/me", middleware.Authenticate(jwtS, constant.EnumRoleUser), userC.DeleteSelfUser)
		userRoutes.POST("", userC.Authenticate)
		userRoutes.POST("/point", middleware.Authenticate(jwtS, constant.EnumRoleUser), userC.AddPoint)
		userRoutes.PATCH("/avatar", middleware.Authenticate(jwtS, constant.EnumRoleUser), userC.ChangeAvatar)
		userRoutes.DELETE("/avatar", middleware.Authenticate(jwtS, constant.EnumRoleUser), userC.DeleteAvatar)
	}
}
