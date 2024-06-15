package router

import (
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func ProfilingRouter(router *gin.Engine, profilingC controller.ProfilingController, jwtS service.JWTService) {
	profilingRoutes := router.Group("/api/v1/profilings")
	{
		profilingRoutes.POST("", middleware.Authenticate(jwtS, constant.EnumRoleUser), profilingC.CreateNewProfiling)
		profilingRoutes.GET("", middleware.Authenticate(jwtS, constant.EnumRoleUser), profilingC.GetUserLast7DaysProfilings)
	}
}
