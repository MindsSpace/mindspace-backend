package service

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type roomService struct {
	roomRepository repository.RoomRepository
	chatRepository repository.ChatRepository
}

type RoomService interface {
	GetRoomAndChatsByID(ctx context.Context, id string) (dto.RoomResponse, error)
	CreateNewRoom(ctx context.Context, ud dto.RoomCreateRequest) (dto.RoomResponse, error)
	GetAllUserRooms(ctx context.Context, userID string) ([]dto.RoomResponse, error)
	DeleteRoomByID(ctx context.Context, id string) error
}

func NewRoomService(roomR repository.RoomRepository, chatR repository.ChatRepository) RoomService {
	return &roomService{roomRepository: roomR, chatRepository: chatR}
}

func (us *roomService) GetRoomAndChatsByID(ctx context.Context, id string) (dto.RoomResponse, error) {
	room, err := us.roomRepository.GetRoomAndChatsByID(ctx, nil, id)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	return dto.RoomResponse{
		ID:        room.ID.String(),
		Name:      room.Name,
		UserID:    room.UserID,
		Chats:     room.Chats,
		CreatedAt: room.CreatedAt.String(),
	}, nil
}

func (us *roomService) CreateNewRoom(ctx context.Context, ud dto.RoomCreateRequest) (dto.RoomResponse, error) {
	db, err := us.roomRepository.TxRepository().BeginTx(ctx)
	if err != nil {
		return dto.RoomResponse{}, err
	}
	defer us.roomRepository.TxRepository().CommitOrRollbackTx(ctx, db, nil)

	room := entity.Room{
		Name:   "Chatroom-" + uuid.New().String(),
		UserID: ud.UserID,
	}

	// create new room
	newRoom, err := us.roomRepository.CreateNewRoom(ctx, db, room)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	chat := entity.Chat{
		Content: ud.Greeting,
		IsUser:  false,
		RoomID:  newRoom.ID.String(),
	}

	// create first chat
	firstChat, err := us.chatRepository.CreateNewChat(ctx, db, chat)
	if err != nil {
		return dto.RoomResponse{}, err
	}

	return dto.RoomResponse{
		ID:        newRoom.ID.String(),
		Name:      newRoom.Name,
		UserID:    newRoom.UserID,
		Chats:     []entity.Chat{firstChat},
		CreatedAt: newRoom.CreatedAt.String(),
	}, nil
}

func (us *roomService) GetAllUserRooms(ctx context.Context, userID string) (roomResp []dto.RoomResponse, err error) {
	rooms, err := us.roomRepository.GetAllUserRooms(ctx, nil, userID)
	if err != nil {
		return []dto.RoomResponse{}, err
	}

	for _, room := range rooms {
		roomResp = append(roomResp, dto.RoomResponse{
			ID:        room.ID.String(),
			Name:      room.Name,
			CreatedAt: room.CreatedAt.String(),
		})
	}

	return roomResp, nil
}

func (us *roomService) DeleteRoomByID(ctx context.Context, id string) error {
	roomCheck, err := us.roomRepository.GetRoomByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(roomCheck, entity.Room{}) {
		return errs.ErrRoomNotFound
	}

	err = us.roomRepository.DeleteRoomByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
