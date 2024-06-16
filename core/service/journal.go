package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-clean-starter/common/util"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type journalService struct {
	journalRepository repository.JournalRepository
	userService       UserService
}

type JournalService interface {
	GetJournalByID(ctx context.Context, id string) (dto.JournalResponse, error)
	CreateNewJournal(ctx context.Context, ud dto.JournalCreateRequest) (dto.JournalResponse, error)
	GetAllUserJournals(ctx context.Context, userID string) ([]dto.JournalResponse, error)
	DeleteJournalByID(ctx context.Context, id string) error
}

func NewJournalService(journalR repository.JournalRepository, userS UserService) JournalService {
	return &journalService{journalRepository: journalR, userService: userS}
}

func (us *journalService) GetJournalByID(ctx context.Context, id string) (dto.JournalResponse, error) {
	journal, err := us.journalRepository.GetJournalByID(ctx, nil, id)
	if err != nil {
		return dto.JournalResponse{}, err
	}

	return dto.JournalResponse{
		ID:        journal.ID.String(),
		Content:   journal.Content,
		Image:     journal.Image,
		UserID:    journal.UserID,
		CreatedAt: journal.CreatedAt.String(),
	}, nil
}

func (us *journalService) CreateNewJournal(ctx context.Context, ud dto.JournalCreateRequest) (dto.JournalResponse, error) {
	db, err := us.journalRepository.TxRepository().BeginTx(ctx)
	if err != nil {
		return dto.JournalResponse{}, err
	}
	defer us.journalRepository.TxRepository().CommitOrRollbackTx(ctx, db, nil)

	journal := entity.Journal{
		Content: ud.Content,
		UserID:  ud.UserID,
	}

	if ud.Image != nil {
		imgID := uuid.New()
		journal.Image = fmt.Sprintf("journal_image/%v", imgID)
		if err := util.UploadFile(ud.Image, journal.Image); err != nil {
			return dto.JournalResponse{}, err
		}
	}

	// create new journal
	newJournal, err := us.journalRepository.CreateNewJournal(ctx, db, journal)
	if err != nil {
		return dto.JournalResponse{}, err
	}

	_, err = us.userService.AddPoint(ctx, ud.UserID, 3)
	if err != nil {
		return dto.JournalResponse{}, err
	}

	return dto.JournalResponse{
		ID:        newJournal.ID.String(),
		Content:   newJournal.Content,
		Image:     newJournal.Image,
		UserID:    newJournal.UserID,
		CreatedAt: newJournal.CreatedAt.String(),
	}, nil
}

func (us *journalService) GetAllUserJournals(ctx context.Context, userID string) (journalResp []dto.JournalResponse, err error) {
	journals, err := us.journalRepository.GetAllUserJournals(ctx, nil, userID)
	if err != nil {
		return []dto.JournalResponse{}, err
	}

	for _, journal := range journals {
		journalResp = append(journalResp, dto.JournalResponse{
			ID:        journal.ID.String(),
			Content:   journal.Content,
			Image:     journal.Image,
			CreatedAt: journal.CreatedAt.String(),
		})
	}

	return journalResp, nil
}

func (us *journalService) DeleteJournalByID(ctx context.Context, id string) error {
	journalCheck, err := us.journalRepository.GetJournalByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(journalCheck, entity.Journal{}) {
		return errs.ErrJournalNotFound
	}

	err = us.journalRepository.DeleteJournalByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}
