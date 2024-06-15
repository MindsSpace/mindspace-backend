package controller

import (
	"net/http"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

type profilingController struct {
	profilingService service.ProfilingService
}

type ProfilingController interface {
	CreateNewProfiling(ctx *gin.Context)
	GetUserLast7DaysProfilings(ctx *gin.Context)
}

func NewProfilingController(profilingS service.ProfilingService) ProfilingController {
	return &profilingController{
		profilingService: profilingS,
	}
}

func (rc *profilingController) CreateNewProfiling(ctx *gin.Context) {
	var profilingDTO dto.ProfilingCreateRequest
	err := ctx.ShouldBind(&profilingDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgProfilingCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	userID := ctx.MustGet("ID").(string)
	profilingDTO.UserID = userID

	newProfiling, err := rc.profilingService.CreateNewProfiling(ctx, profilingDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgProfilingCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MsgProfilingCreateSuccess,
		http.StatusCreated, newProfiling,
	))
}

func (rc *profilingController) GetUserLast7DaysProfilings(ctx *gin.Context) {
	userID := ctx.MustGet("ID").(string)

	profilings, err := rc.profilingService.GetUserLast7DaysProfilings(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgProfilingsFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgProfilingsFetchSuccess,
		http.StatusOK, profilings,
	))
}
