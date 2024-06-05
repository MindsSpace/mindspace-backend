package service

import (
	"context"
	"reflect"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type chatService struct {
	chatRepository repository.ChatRepository
}

type ChatService interface {
	GetChatByID(ctx context.Context, id string) (dto.ChatResponse, error)
	CreateNewChat(ctx context.Context, ud dto.ChatCreateRequest) (dto.ChatResponse, error)
	DeleteChatByID(ctx context.Context, id string) error
}

func NewChatService(chatR repository.ChatRepository) ChatService {
	return &chatService{chatRepository: chatR}
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
	chat := entity.Chat{
		Content: cd.Content,
		IsUser:  true,
		RoomID:  cd.RoomID,
	}

	// create new chat
	newChat, err := us.chatRepository.CreateNewChat(ctx, nil, chat)
	if err != nil {
		return dto.ChatResponse{}, err
	}

	return dto.ChatResponse{
		ID:        newChat.ID.String(),
		Content:   newChat.Content,
		IsUser:    newChat.IsUser,
		RoomID:    newChat.RoomID,
		CreatedAt: newChat.CreatedAt.String(),
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
