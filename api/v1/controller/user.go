package controller

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

type UserController interface {
	Authenticate(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	UpdateUserByID(ctx *gin.Context)
	DeleteSelfUser(ctx *gin.Context)
	DeleteUserByID(ctx *gin.Context)
	AddPoint(ctx *gin.Context)
	ChangeAvatar(ctx *gin.Context)
	DeleteAvatar(ctx *gin.Context)
}

func NewUserController(userS service.UserService, jwtS service.JWTService) UserController {
	return &userController{
		userService: userS,
		jwtService:  jwtS,
	}
}

func (uc *userController) Authenticate(ctx *gin.Context) {
	var userDTO dto.UserAuthRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAuthenticateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.AuthenticateUser(ctx, userDTO)
	if err != nil {
		if errors.Is(err, errs.ErrPasswordWrong) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
				messages.MsgUserWrongPassword,
				err.Error(), http.StatusBadRequest,
			))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAuthenticateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, constant.EnumRoleUser)
	authResp := base.CreateAuthResponse(token, constant.EnumRoleUser, user.IsProfiled)
	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserAuthenticateSuccess,
		http.StatusOK, authResp,
	))
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	var req base.GetsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	users, pageMeta, err := uc.userService.GetAllUsers(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUsersFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	if reflect.DeepEqual(pageMeta, base.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
			messages.MsgUsersFetchSuccess,
			http.StatusOK, users,
		))
	} else {
		ctx.JSON(http.StatusOK, base.CreatePaginatedResponse(
			messages.MsgUsersFetchSuccess,
			http.StatusOK, users, pageMeta,
		))
	}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrID, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserFetchSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")

	var userDTO dto.UserUpdateRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	user, err := uc.userService.UpdateUserByID(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (uc *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := uc.userService.DeleteUserByID(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (uc *userController) AddPoint(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var userDTO dto.UserAddPointRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAvatarUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res, err := uc.userService.AddPoint(ctx, id, userDTO.Point)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAvatarUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserAvatarUpdateSuccess,
		http.StatusOK, res,
	))
}

func (uc *userController) ChangeAvatar(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var userDTO dto.UserChangeAvatarRequest
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAvatarUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	res, err := uc.userService.ChangeAvatar(ctx, userDTO, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAvatarUpdateFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserAvatarUpdateSuccess,
		http.StatusOK, res,
	))
}

func (uc *userController) DeleteAvatar(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	err := uc.userService.DeleteAvatar(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgUserAvatarDeleteFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserAvatarDeleteSuccess,
		http.StatusOK, nil,
	))
}
