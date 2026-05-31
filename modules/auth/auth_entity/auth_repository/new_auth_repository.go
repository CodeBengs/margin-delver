package authrepository

import (
	"context"

	authdto "margin-delver/modules/auth/auth_dto"
	authmodel "margin-delver/modules/auth/auth_entity/auth_model"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	Login(ctx context.Context, request *authdto.LoginRequest) (*authdto.LoginResponse, error)
	Create(ctx context.Context, user *authmodel.User) error
	Count(ctx context.Context) (int64, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) AuthRepositoryInterface {
	return &AuthRepository{
		db: db,
	}
}
