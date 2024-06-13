package repository

import (
	"context"
	"errors"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"

	"gorm.io/gorm"
)

type roomRepository struct {
	txr *txRepository
}

type RoomRepository interface {
	// tx
	TxRepository() *txRepository

	// functional
	CreateNewRoom(ctx context.Context, tx *gorm.DB, room entity.Room) (entity.Room, error)
	GetRoomByID(ctx context.Context, tx *gorm.DB, id string) (entity.Room, error)
	GetRoomAndChatsByID(ctx context.Context, tx *gorm.DB, id string) (entity.Room, error)
	GetAllUserRooms(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Room, error)
	UpdateRoomName(ctx context.Context, tx *gorm.DB, id string, new string) error
	DeleteRoomByID(ctx context.Context, tx *gorm.DB, id string) error
}

func NewRoomRepository(txr *txRepository) *roomRepository {
	return &roomRepository{txr: txr}
}

func (rr *roomRepository) TxRepository() *txRepository {
	return rr.txr
}

func (rr *roomRepository) CreateNewRoom(ctx context.Context, tx *gorm.DB, room entity.Room) (entity.Room, error) {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Create(&room)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&room).Error
	}

	if err != nil {
		return entity.Room{}, err
	}
	return room, nil
}

func (rr *roomRepository) GetRoomByID(ctx context.Context, tx *gorm.DB, id string) (entity.Room, error) {
	var err error
	var room entity.Room
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Where("id = $1", id).Take(&room)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&room).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return room, err
	}
	return room, nil
}

func (rr *roomRepository) GetRoomAndChatsByID(ctx context.Context, tx *gorm.DB, id string) (entity.Room, error) {
	var err error
	var room entity.Room
	if tx == nil {
		err = rr.txr.DB().WithContext(ctx).Debug().Where("id = $1", id).Preload("Chats", func(db *gorm.DB) *gorm.DB {
			return db.Order("chats.created_at ASC")
		}).Take(&room).Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Preload("Chats", func(db *gorm.DB) *gorm.DB {
			return db.Order("chats.created_at ASC")
		}).Take(&room).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return room, err
	}
	return room, nil
}

func (rr *roomRepository) GetAllUserRooms(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Room, error) {
	var err error
	var rooms []entity.Room

	if tx == nil {
		tx = rr.txr.DB()
	}

	err = tx.WithContext(ctx).Debug().Where("user_id = ?", userID).Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (rr *roomRepository) UpdateRoomName(ctx context.Context, tx *gorm.DB, id string, new string) error {
	if tx == nil {
		tx = rr.txr.DB()
	}

	if err := tx.WithContext(ctx).Debug().Model(&entity.Room{}).Where("id = ?", id).Update("name", new).Error; err != nil {
		return err
	}
	return nil
}

func (rr *roomRepository) DeleteRoomByID(ctx context.Context, tx *gorm.DB, id string) error {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Delete(&entity.Room{}, &id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.Room{}, &id).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return err
	}
	return nil
}
