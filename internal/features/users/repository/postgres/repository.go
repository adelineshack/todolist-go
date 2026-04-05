package users_postgres_repository

import core_postgres_pool "github.com/adelineshack/todolist-go/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	if pool == nil {
		panic("users repository: nil pool")
	}

	return &UsersRepository{
		pool: pool,
	}
}
