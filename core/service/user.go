package service

import (
	"context"
	"reflect"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/util"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	VerifyLogin(ctx context.Context, username string, password string) bool
	CreateNewUser(ctx context.Context, ud dto.UserAuthRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context, req base.GetsRequest) ([]dto.UserResponse, base.PaginationResponse, error)
	GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error)
	UpdateUserByID(ctx context.Context, ud dto.UserUpdateRequest, id string) (dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, id string) error
}

func NewUserService(userR repository.UserRepository) UserService {
	return &userService{userRepository: userR}
}

func (us *userService) VerifyLogin(ctx context.Context, username string, password string) bool {
	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrUsername, username)
	if err != nil {
		return false
	}
	passwordCheck, err := util.PasswordCompare(userCheck.Password, []byte(password))
	if err != nil {
		return false
	}

	if userCheck.Username == username && passwordCheck {
		return true
	}
	return false
}

func (us *userService) CreateNewUser(ctx context.Context, ud dto.UserAuthRequest) (dto.UserResponse, error) {
	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrUsername, ud.Username)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		return dto.UserResponse{}, errs.ErrUsernameAlreadyExists
	}

	user := entity.User{
		Username: ud.Username,
		Password: ud.Password,
		Level:    1,
		Point:    0,
	}

	// create new user
	newUser, err := us.userRepository.CreateNewUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
		Level:    newUser.Level,
		Point:    newUser.Point,
	}, nil
}

func (us *userService) GetAllUsers(ctx context.Context, req base.GetsRequest) (
	userResp []dto.UserResponse, pageResp base.PaginationResponse, err error) {
	if req.Limit < 0 {
		req.Limit = 0
	}

	if req.Page < 0 {
		req.Page = 0
	}

	if req.Sort != "" && req.Sort[0] == '-' {
		req.Sort = req.Sort[1:] + " DESC"
	}

	users, lastPage, total, err := us.userRepository.GetAllUsers(ctx, nil, req)
	if err != nil {
		return []dto.UserResponse{}, base.PaginationResponse{}, err
	}

	for _, user := range users {
		userResp = append(userResp, dto.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Level:    user.Level,
			Point:    user.Point,
		})
	}

	if req.Limit == 0 {
		return userResp, base.PaginationResponse{}, nil
	}

	pageResp = base.PaginationResponse{
		Page:     int64(req.Page),
		Limit:    int64(req.Limit),
		LastPage: lastPage,
		Total:    total,
	}
	return userResp, pageResp, nil
}

func (us *userService) GetUserByPrimaryKey(ctx context.Context, key string, val string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, key, val)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Level:    user.Level,
		Point:    user.Point,
	}, nil
}

func (us *userService) UpdateUserByID(ctx context.Context,
	ud dto.UserUpdateRequest, id string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, errs.ErrUserNotFound
	}

	if ud.Username != "" && ud.Username != user.Username {
		us, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrUsername, ud.Username)
		if err != nil {
			return dto.UserResponse{}, err
		}

		if !(reflect.DeepEqual(us, entity.User{})) {
			return dto.UserResponse{}, errs.ErrUsernameAlreadyExists
		}
	}

	userEdit := entity.User{
		ID:       user.ID,
		Username: ud.Username,
		Level:    ud.Level,
		Password: ud.Password,
		Point:    ud.Point,
	}

	edited, err := us.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if edited.Username == "" {
		edited.Username = user.Username
	}
	if edited.Point == 0 {
		edited.Point = user.Point
	}

	return dto.UserResponse{
		ID:       edited.ID.String(),
		Username: edited.Username,
		Level:    edited.Level,
		Point:    user.Point,
	}, nil
}

func (us *userService) DeleteUserByID(ctx context.Context, id string) error {
	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(userCheck, entity.User{}) {
		return errs.ErrUserNotFound
	}

	err = us.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
