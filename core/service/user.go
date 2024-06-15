package service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/common/util"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type userService struct {
	userRepository      repository.UserRepository
	profilingRepository repository.ProfilingRepository
}

type UserService interface {
	AuthenticateUser(ctx context.Context, ud dto.UserAuthRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context, req base.GetsRequest) ([]dto.UserResponse, base.PaginationResponse, error)
	GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error)
	UpdateUserByID(ctx context.Context, ud dto.UserUpdateRequest, id string) (dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, id string) error
	AddPoint(ctx context.Context, userID string, point int) (dto.UserResponse, error)
	ChangeAvatar(ctx context.Context, req dto.UserChangeAvatarRequest, userID string) (dto.UserResponse, error)
	DeleteAvatar(ctx context.Context, userID string) error
}

func NewUserService(userR repository.UserRepository, profilingR repository.ProfilingRepository) UserService {
	return &userService{userRepository: userR, profilingRepository: profilingR}
}

func (us *userService) AuthenticateUser(ctx context.Context, ud dto.UserAuthRequest) (dto.UserResponse, error) {
	db, err := us.userRepository.TxRepository().BeginTx(ctx)
	if err != nil {
		return dto.UserResponse{}, err
	}
	defer us.userRepository.TxRepository().CommitOrRollbackTx(ctx, db, nil)

	userCheck, err := us.userRepository.GetUserByPrimaryKey(ctx, db, constant.DBAttrUsername, ud.Username)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// check if exist
	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		passwordCheck, err := util.PasswordCompare(userCheck.Password, []byte(ud.Password))
		if err != nil {
			return dto.UserResponse{}, err
		}

		if userCheck.Username == ud.Username && passwordCheck {
			lastProfiling, err := us.profilingRepository.GetUserLatestProfiling(ctx, db, userCheck.ID.String())
			if err != nil {
				return dto.UserResponse{}, err
			}

			isProfiled := false
			if lastProfiling.CreatedAt.Day() == time.Now().Day() {
				isProfiled = true
			}

			return dto.UserResponse{
				ID:         userCheck.ID.String(),
				Username:   userCheck.Username,
				Level:      userCheck.Level,
				Point:      userCheck.Point,
				IsProfiled: &isProfiled,
			}, nil
		}
		return dto.UserResponse{}, errs.ErrPasswordWrong
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

	lastProfiling, err := us.profilingRepository.GetUserLatestProfiling(ctx, nil, user.ID.String())
	if err != nil {
		return dto.UserResponse{}, err
	}

	isProfiled := false
	if lastProfiling.CreatedAt.Day() == time.Now().Day() {
		isProfiled = true
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Username:   user.Username,
		Level:      user.Level,
		Point:      user.Point,
		IsProfiled: &isProfiled,
		Avatar:     *user.Avatar,
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

func (us *userService) AddPoint(ctx context.Context, userID string, point int) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, errs.ErrUserNotFound
	}

	pointNeeded := ((constant.DefaultLevelPointMultiplier * user.Level) - user.Point)
	for point >= pointNeeded {
		point -= pointNeeded
		user.Level++
		pointNeeded = constant.DefaultLevelPointMultiplier * user.Level
	}
	user.Point = point

	userUpdate, err := us.userRepository.UpdateUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       userUpdate.ID.String(),
		Username: userUpdate.Username,
		Level:    userUpdate.Level,
		Point:    userUpdate.Point,
		Avatar:   *userUpdate.Avatar,
	}, nil
}

func (us *userService) ChangeAvatar(ctx context.Context, req dto.UserChangeAvatarRequest, userID string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, errs.ErrUserNotFound
	}

	if *user.Avatar != "" {
		if err := util.DeleteFile(*user.Avatar); err != nil {
			return dto.UserResponse{}, err
		}
	}

	picID := uuid.New()
	picPath := fmt.Sprintf("user_avatar/%v", picID)

	userEdit := entity.User{
		ID:     user.ID,
		Avatar: &picPath,
	}

	if err := util.UploadFile(req.Avatar, picPath); err != nil {
		return dto.UserResponse{}, err
	}

	userUpdate, err := us.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:     userUpdate.ID.String(),
		Avatar: *userUpdate.Avatar,
	}, nil
}

func (us *userService) DeleteAvatar(ctx context.Context, userID string) error {
	user, err := us.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, userID)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return errs.ErrUserNotFound
	}

	if *user.Avatar == "" {
		return errs.ErrUserNoAvatar
	}

	if err := util.DeleteFile(*user.Avatar); err != nil {
		return err
	}

	emptyStr := ""
	userEdit := entity.User{
		ID:     user.ID,
		Avatar: &emptyStr,
	}

	_, err = us.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return err
	}

	return nil
}
