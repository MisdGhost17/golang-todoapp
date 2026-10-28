package auth_postgres_repository

import core_postgres_pool "github.com/MisdGhost17/golang-todoapp/internal/core/repository/postgres/pool"

type AuthRepository struct {
	pool core_postgres_pool.Pool
}

func NewAuthRepository(pool core_postgres_pool.Pool) *AuthRepository {
	return &AuthRepository{
		pool: pool,
	}
}
