package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
)

type profilingService struct {
	profilingRepository repository.ProfilingRepository
	roomService         RoomService
}

type ProfilingService interface {
	CreateNewProfiling(ctx context.Context, ud dto.ProfilingCreateRequest) (dto.ProfilingResponse, error)
	GetUserLast7DaysProfilings(ctx context.Context, userID string) ([]dto.ProfilingResponse, error)
}

func NewProfilingService(profilingR repository.ProfilingRepository, roomS RoomService) ProfilingService {
	return &profilingService{profilingRepository: profilingR, roomService: roomS}
}

func (us *profilingService) CreateNewProfiling(ctx context.Context, ud dto.ProfilingCreateRequest) (dto.ProfilingResponse, error) {
	db, err := us.profilingRepository.TxRepository().BeginTx(ctx)
	if err != nil {
		return dto.ProfilingResponse{}, err
	}
	defer us.profilingRepository.TxRepository().CommitOrRollbackTx(ctx, db, nil)

	lastProfiling, err := us.profilingRepository.GetUserLatestProfiling(ctx, db, ud.UserID)
	if err != nil {
		return dto.ProfilingResponse{}, err
	}

	if lastProfiling.CreatedAt.Day() == time.Now().Day() {
		return dto.ProfilingResponse{}, errors.ErrProfilingFilledToday
	}

	// create room for profiling
	newRoom, err := us.roomService.CreateNewRoom(ctx, dto.RoomCreateRequest{
		Greeting: constant.GreetingMessages[rand.Intn(len(constant.GreetingMessages))],
		UserID:   ud.UserID,
	})
	if err != nil {
		return dto.ProfilingResponse{}, err
	}

	profiling := entity.Profiling{
		Mood:       ud.Mood,
		Problems:   ud.Problems,
		Approaches: ud.Approaches,
		UserID:     ud.UserID,
		RoomID:     newRoom.ID,
	}

	// create new profiling
	newProfiling, err := us.profilingRepository.CreateNewProfiling(ctx, nil, profiling)
	if err != nil {
		return dto.ProfilingResponse{}, err
	}

	return dto.ProfilingResponse{
		ID:         newProfiling.ID.String(),
		Mood:       newProfiling.Mood,
		Problems:   newProfiling.Problems,
		Approaches: newProfiling.Approaches,
		UserID:     newProfiling.UserID,
		RoomID:     newRoom.ID,
		IsFilled:   true,
		CreatedAt:  newProfiling.CreatedAt.String(),
	}, nil
}

func (us *profilingService) GetUserLast7DaysProfilings(ctx context.Context, userID string) (profilingResp []dto.ProfilingResponse, err error) {
	profilings, past7Day, err := us.profilingRepository.GetUserLast7DaysProfilings(ctx, nil, userID)
	if err != nil {
		return []dto.ProfilingResponse{}, err
	}

	idx := 0
	for i := 0; i < 7; i++ {
		curTime := past7Day.AddDate(0, 0, i)
		profiling := dto.ProfilingResponse{
			IsFilled:  false,
			CreatedAt: curTime.String(),
		}

		if idx < len(profilings) && profilings[idx].CreatedAt.Day() == curTime.Day() {
			profiling = dto.ProfilingResponse{
				ID:         profilings[idx].ID.String(),
				Mood:       profilings[idx].Mood,
				Problems:   profilings[idx].Problems,
				Approaches: profilings[idx].Approaches,
				IsFilled:   true,
				UserID:     profilings[idx].UserID,
				RoomID:     profilings[idx].RoomID,
				CreatedAt:  profilings[idx].CreatedAt.String(),
			}
			idx++
		}
		profilingResp = append(profilingResp, profiling)
	}

	return profilingResp, nil
}
