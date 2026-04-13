package domain

import (
	"fmt"
	"time"

	core_errors "github.com/adelineshack/todolist-go/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	userAuthorId int,
) Task {
	return Task{
		ID:      id,
		Version: version,

		Title:       title,
		Description: description,
		Completed:   completed,
		CreatedAt:   createdAt,
		CompletedAt: completedAt,

		AuthorUserID: userAuthorId,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	userAuthorId int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		userAuthorId,
	)
}

func (t *Task) Validate() error {

	titleLen := len([]rune(t.Title))

	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))

		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"CompletedAt can't be `nil` if `Completed` == 'true': %w",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` can't be earlier than `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` can't have value when `Completed` == 'true': %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
