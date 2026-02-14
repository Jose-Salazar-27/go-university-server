package persistence

import (
	"context"
	"database/sql"

	"github.com/Jose-Salazar-27/go-university-server/internal/auth/domain"
	"github.com/Jose-Salazar-27/go-university-server/internal/shared/kernel/infra/db"
)

type postgresUserRepository struct {
	pool *sql.DB
}

func NewUserRepository(db *sql.DB) *postgresUserRepository {
	return &postgresUserRepository{db}
}

func (r postgresUserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (
		id,
		email, 
		password_hash,
		first_name,
		last_name,
		user_type,
		avatar_url,
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	if _, err := r.pool.ExecContext(ctx, query); err != nil {
		if ok, pgerr := db.IsPgError(err); ok {
			return db.ExchangePGError(pgerr)
		}
		return err
	}
	return nil
}
