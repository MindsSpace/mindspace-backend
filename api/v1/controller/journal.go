package controller

import (
	"net/http"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

type journalController struct {
	journalService service.JournalService
}

type JournalController interface {
	GetJournalByID(ctx *gin.Context)
	CreateNewJournal(ctx *gin.Context)
	GetAllUserJournals(ctx *gin.Context)
	DeleteJournalByID(ctx *gin.Context)
}

func NewJournalController(journalS service.JournalService) JournalController {
	return &journalController{
		journalService: journalS,
	}
}

func (rc *journalController) GetJournalByID(ctx *gin.Context) {
	id := ctx.Param("journal_id")

	journal, err := rc.journalService.GetJournalByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgJournalFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgJournalFetchSuccess,
		http.StatusOK, journal,
	))
}

func (rc *journalController) CreateNewJournal(ctx *gin.Context) {
	userID := ctx.MustGet("ID").(string)

	var journalDTO dto.JournalCreateRequest
	err := ctx.ShouldBind(&journalDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgJournalCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}
	journalDTO.UserID = userID

	newJournal, err := rc.journalService.CreateNewJournal(ctx, journalDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgJournalCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MsgJournalCreateSuccess,
		http.StatusCreated, newJournal,
	))
}

func (rc *journalController) GetAllUserJournals(ctx *gin.Context) {
	userID := ctx.MustGet("ID").(string)

	journals, err := rc.journalService.GetAllUserJournals(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgJournalsFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgJournalsFetchSuccess,
		http.StatusOK, journals,
	))
}

func (rc *journalController) DeleteJournalByID(ctx *gin.Context) {
	id := ctx.Param("journal_id")

	err := rc.journalService.DeleteJournalByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgJournalDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgJournalDeleteSuccess,
		http.StatusOK, nil,
	))
}
