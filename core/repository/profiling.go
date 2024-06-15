package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"

	"gorm.io/gorm"
)

type profilingRepository struct {
	txr *txRepository
}

type ProfilingRepository interface {
	// tx
	TxRepository() *txRepository

	// functional
	CreateNewProfiling(ctx context.Context, tx *gorm.DB, profiling entity.Profiling) (entity.Profiling, error)
	GetProfilingByID(ctx context.Context, tx *gorm.DB, id string) (entity.Profiling, error)
	GetUserLast7DaysProfilings(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Profiling, time.Time, error)
	GetUserLatestProfiling(ctx context.Context, tx *gorm.DB, userID string) (entity.Profiling, error)
	DeleteProfilingByID(ctx context.Context, tx *gorm.DB, id string) error
}

func NewProfilingRepository(txr *txRepository) *profilingRepository {
	return &profilingRepository{txr: txr}
}

func (rr *profilingRepository) TxRepository() *txRepository {
	return rr.txr
}

func (rr *profilingRepository) CreateNewProfiling(ctx context.Context, tx *gorm.DB, profiling entity.Profiling) (entity.Profiling, error) {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Create(&profiling)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Create(&profiling).Error
	}

	if err != nil {
		return entity.Profiling{}, err
	}
	return profiling, nil
}

func (rr *profilingRepository) GetProfilingByID(ctx context.Context, tx *gorm.DB, id string) (entity.Profiling, error) {
	var err error
	var profiling entity.Profiling
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Where("id = $1", id).Take(&profiling)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Where("id = $1", id).Take(&profiling).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return profiling, err
	}
	return profiling, nil
}

func (rr *profilingRepository) GetUserLast7DaysProfilings(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Profiling, time.Time, error) {
	var profilings []entity.Profiling

	if tx == nil {
		tx = rr.txr.DB()
	}

	past7Day := time.Now().AddDate(0, 0, -6)
	if err := tx.WithContext(ctx).Debug().Where("created_at >= ?", past7Day).
		Where("user_id = ?", userID).Order("created_at ASC").Find(&profilings).Error; err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return profilings, time.Time{}, err
	}
	return profilings, past7Day, nil
}

func (rr *profilingRepository) GetUserLatestProfiling(ctx context.Context, tx *gorm.DB, userID string) (entity.Profiling, error) {
	var profiling entity.Profiling

	if tx == nil {
		tx = rr.txr.DB()
	}

	err := tx.WithContext(ctx).Debug().Where("user_id = $1", userID).Order("created_at DESC").Take(&profiling).Error
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return profiling, err
	}
	return profiling, nil
}

func (rr *profilingRepository) DeleteProfilingByID(ctx context.Context, tx *gorm.DB, id string) error {
	var err error
	if tx == nil {
		tx = rr.txr.DB().WithContext(ctx).Debug().Delete(&entity.Profiling{}, &id)
		err = tx.Error
	} else {
		err = tx.WithContext(ctx).Debug().Delete(&entity.Profiling{}, &id).Error
	}

	if err != nil {
		return err
	}
	return nil
}
