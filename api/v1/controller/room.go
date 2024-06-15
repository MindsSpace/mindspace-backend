package controller

import (
	"math/rand"
	"net/http"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

type roomController struct {
	roomService service.RoomService
}

type RoomController interface {
	GetRoomAndChatsByID(ctx *gin.Context)
	CreateNewRoom(ctx *gin.Context)
	GetAllUserRooms(ctx *gin.Context)
	DeleteRoomByID(ctx *gin.Context)
}

func NewRoomController(roomS service.RoomService) RoomController {
	return &roomController{
		roomService: roomS,
	}
}

func (rc *roomController) GetRoomAndChatsByID(ctx *gin.Context) {
	id := ctx.Param("room_id")

	room, err := rc.roomService.GetRoomAndChatsByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgRoomFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgRoomFetchSuccess,
		http.StatusOK, room,
	))
}

func (rc *roomController) CreateNewRoom(ctx *gin.Context) {
	userID := ctx.MustGet("ID").(string)

	newRoom, err := rc.roomService.CreateNewRoom(ctx, dto.RoomCreateRequest{
		Greeting: constant.GreetingMessages[rand.Intn(len(constant.GreetingMessages))],
		UserID:   userID,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgRoomCreateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MsgRoomCreateSuccess,
		http.StatusCreated, newRoom,
	))
}

func (rc *roomController) GetAllUserRooms(ctx *gin.Context) {
	userID := ctx.MustGet("ID").(string)

	rooms, err := rc.roomService.GetAllUserRooms(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgRoomsFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgRoomsFetchSuccess,
		http.StatusOK, rooms,
	))
}

func (rc *roomController) DeleteRoomByID(ctx *gin.Context) {
	id := ctx.Param("room_id")

	err := rc.roomService.DeleteRoomByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgRoomDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgRoomDeleteSuccess,
		http.StatusOK, nil,
	))
}
