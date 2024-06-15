package router

import (
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func JournalRouter(router *gin.Engine, journalC controller.JournalController, jwtS service.JWTService) {
	journalRoutes := router.Group("/api/v1/journals")
	{
		journalRoutes.GET("", middleware.Authenticate(jwtS, constant.EnumRoleUser), journalC.GetAllUserJournals)
		journalRoutes.GET("/:journal_id", middleware.Authenticate(jwtS, constant.EnumRoleUser), journalC.GetJournalByID)
		journalRoutes.POST("", middleware.Authenticate(jwtS, constant.EnumRoleUser), journalC.CreateNewJournal)
		journalRoutes.DELETE("/:journal_id", middleware.Authenticate(jwtS, constant.EnumRoleUser), journalC.DeleteJournalByID)
	}
}
