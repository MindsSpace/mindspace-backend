package repository

import (
	"context"
	"errors"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"

	"gorm.io/gorm"
)

type chatRepository struct {
	txr *txRepository
}

type ChatRepository interface {
	// tx
	TxRepository() *txRepository

	// functional
	CreateNewChat(ctx context.Context, tx *gorm.DB, chat entity.Chat) (entity.Chat, error)
	GetChatByID(ctx context.Context, tx *gorm.DB, id string) (entity.Chat, error)
	DeleteChatByID(ctx context.Context, tx *gorm.DB, id string) error
}

func NewChatRepository(txr *txRepository) *chatRepository {
	return &chatRepository{txr: txr}
}

func (rr *chatRepository) TxRepository() *txRepository {
	return rr.txr
}

func (rr *chatRepository) CreateNewChat(ctx context.Context, tx *gorm.DB, chat entity.Chat) (entity.Chat, error) {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Create(&chat)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&chat).Error
	}

	if err != nil {
		return entity.Chat{}, err
	}
	return chat, nil
}

func (rr *chatRepository) GetChatByID(ctx context.Context, tx *gorm.DB, id string) (entity.Chat, error) {
	var err error
	var chat entity.Chat
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Where("id = $1", id).Take(&chat)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&chat).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return chat, err
	}
	return chat, nil
}

func (rr *chatRepository) DeleteChatByID(ctx context.Context, tx *gorm.DB, id string) error {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Delete(&entity.Chat{}, &id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.Chat{}, &id).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return err
	}
	return nil
}
