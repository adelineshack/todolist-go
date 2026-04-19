package statistics_service

import (
	"context"
	"time"

	"github.com/adelineshack/todolist-go/internal/core/domain"
)

type StasticticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewTasksService(
	statisticsRepository StatisticsRepository,
) *StasticticsService {
	return &StasticticsService{
		statisticsRepository: statisticsRepository,
	}
}
