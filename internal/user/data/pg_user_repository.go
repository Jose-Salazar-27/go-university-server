package data

import (
	"context"
	"errors"

	"github.com/Jose-Salazar-27/go-university-server/internal/user/entity"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *postgresUserRepository {
	return &postgresUserRepository{db}
}

func (r postgresUserRepository) Create(ctx context.Context, u *entity.User) error {
	model := fromEntity(u)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("item already exits")
		}
		return err
	}
	return nil
}
