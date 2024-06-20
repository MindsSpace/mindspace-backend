package service

import (
	"context"
	"reflect"

	"github.com/zetsux/gin-gorm-clean-starter/common/util"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type chatService struct {
	chatRepository repository.ChatRepository
	roomRepository repository.RoomRepository
}

type ChatService interface {
	GetChatByID(ctx context.Context, id string) (dto.ChatResponse, error)
	CreateNewChat(ctx context.Context, ud dto.ChatCreateRequest) (dto.ChatResponse, error)
	DeleteChatByID(ctx context.Context, id string) error
}

func NewChatService(chatR repository.ChatRepository, roomR repository.RoomRepository) ChatService {
	return &chatService{chatRepository: chatR, roomRepository: roomR}
}

func (us *chatService) GetChatByID(ctx context.Context, id string) (dto.ChatResponse, error) {
	chat, err := us.chatRepository.GetChatByID(ctx, nil, id)
	if err != nil {
		return dto.ChatResponse{}, err
	}

	return dto.ChatResponse{
		ID:        chat.ID.String(),
		Content:   chat.Content,
		IsUser:    chat.IsUser,
		RoomID:    chat.RoomID,
		CreatedAt: chat.CreatedAt.String(),
	}, nil
}

func (us *chatService) CreateNewChat(ctx context.Context, cd dto.ChatCreateRequest) (dto.ChatResponse, error) {
	db, err := us.chatRepository.TxRepository().BeginTx(ctx)
	if err != nil {
		return dto.ChatResponse{}, err
	}
	defer us.roomRepository.TxRepository().CommitOrRollbackTx(ctx, db, nil)

	checkRoom, err := us.roomRepository.GetRoomAndChatsByID(ctx, db, cd.RoomID)
	if err != nil {
		return dto.ChatResponse{}, err
	} else if len(checkRoom.Chats) <= 1 {
		us.roomRepository.UpdateRoomName(ctx, db, checkRoom.ID.String(), cd.Content)
	}

	chat := entity.Chat{
		Content: cd.Content,
		IsUser:  true,
		RoomID:  cd.RoomID,
	}

	// create new chat
	_, err = us.chatRepository.CreateNewChat(ctx, db, chat)
	if err != nil {
		return dto.ChatResponse{}, err
	}

	chatbotResp, err := util.GetChatbotResponse(cd.Language, cd.Content)
	if err != nil {
		return dto.ChatResponse{}, err
	}

	response := entity.Chat{
		Content: chatbotResp,
		IsUser:  false,
		RoomID:  cd.RoomID,
	}

	// create new chat
	newResponse, err := us.chatRepository.CreateNewChat(ctx, db, response)
	if err != nil {
		return dto.ChatResponse{}, err
	}

	return dto.ChatResponse{
		ID:        newResponse.ID.String(),
		Content:   newResponse.Content,
		IsUser:    newResponse.IsUser,
		RoomID:    newResponse.RoomID,
		CreatedAt: newResponse.CreatedAt.String(),
	}, nil
}

func (us *chatService) DeleteChatByID(ctx context.Context, id string) error {
	chatCheck, err := us.chatRepository.GetChatByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(chatCheck, entity.Chat{}) {
		return errs.ErrChatNotFound
	}

	err = us.chatRepository.DeleteChatByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
