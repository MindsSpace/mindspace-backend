package controller

import (
	"net/http"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

type chatController struct {
	chatService service.ChatService
}

type ChatController interface {
	GetChatByID(ctx *gin.Context)
	CreateNewChat(ctx *gin.Context)
	DeleteChatByID(ctx *gin.Context)
}

func NewChatController(chatS service.ChatService) ChatController {
	return &chatController{
		chatService: chatS,
	}
}

func (rc *chatController) GetChatByID(ctx *gin.Context) {
	id := ctx.Param("chat_id")

	chat, err := rc.chatService.GetChatByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgChatFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgChatFetchSuccess,
		http.StatusOK, chat,
	))
}

func (rc *chatController) CreateNewChat(ctx *gin.Context) {
	var chatDTO dto.ChatCreateRequest
	err := ctx.ShouldBind(&chatDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgChatCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	newChat, err := rc.chatService.CreateNewChat(ctx, chatDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgChatCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MsgChatCreateSuccess,
		http.StatusCreated, newChat,
	))
}

func (rc *chatController) DeleteChatByID(ctx *gin.Context) {
	id := ctx.Param("chat_id")

	err := rc.chatService.DeleteChatByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgChatDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgChatDeleteSuccess,
		http.StatusOK, nil,
	))
}
