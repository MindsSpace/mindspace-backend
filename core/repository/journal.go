package repository

import (
	"context"
	"errors"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"

	"gorm.io/gorm"
)

type journalRepository struct {
	txr *txRepository
}

type JournalRepository interface {
	// tx
	TxRepository() *txRepository

	// functional
	CreateNewJournal(ctx context.Context, tx *gorm.DB, journal entity.Journal) (entity.Journal, error)
	GetJournalByID(ctx context.Context, tx *gorm.DB, id string) (entity.Journal, error)
	GetAllUserJournals(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Journal, error)
	DeleteJournalByID(ctx context.Context, tx *gorm.DB, id string) error
}

func NewJournalRepository(txr *txRepository) *journalRepository {
	return &journalRepository{txr: txr}
}

func (rr *journalRepository) TxRepository() *txRepository {
	return rr.txr
}

func (rr *journalRepository) CreateNewJournal(ctx context.Context, tx *gorm.DB, journal entity.Journal) (entity.Journal, error) {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Create(&journal)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&journal).Error
	}

	if err != nil {
		return entity.Journal{}, err
	}
	return journal, nil
}

func (rr *journalRepository) GetJournalByID(ctx context.Context, tx *gorm.DB, id string) (entity.Journal, error) {
	var err error
	var journal entity.Journal
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Where("id = $1", id).Take(&journal)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&journal).Error
	}

	if err != nil {
		return journal, err
	}
	return journal, nil
}

func (rr *journalRepository) GetAllUserJournals(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Journal, error) {
	var err error
	var journals []entity.Journal

	if tx == nil {
		tx = rr.txr.DB()
	}

	err = tx.WithContext(ctx).Debug().Where("user_id = ?", userID).Find(&journals).Error
	if err != nil {
		return nil, err
	}
	return journals, nil
}

func (rr *journalRepository) DeleteJournalByID(ctx context.Context, tx *gorm.DB, id string) error {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Delete(&entity.Journal{}, &id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.Journal{}, &id).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return err
	}
	return nil
}
