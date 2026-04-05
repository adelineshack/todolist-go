package users_service

import (
	"context"
	"fmt"

	"github.com/adelineshack/todolist-go/internal/core/domain"
	core_logger "github.com/adelineshack/todolist-go/internal/core/logger"
	"go.uber.org/zap"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	log := core_logger.FromContext(ctx)
	log.Debug("inside CreateUser",
		zap.Bool("serviceNil", s == nil),
		zap.Bool("repoNil", s != nil && s.usersRepository == nil),
	)

	user, err := s.usersRepository.CreateUser(ctx, user)

	log.Debug("after repository CreateUser", zap.Error(err))

	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
