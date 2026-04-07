package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/adelineshack/todolist-go/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(
		ctx, r.pool.OpTimeout(),
	)

	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1;
	`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf(
			"exec query: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with id='%d': '%w'",
			id,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
