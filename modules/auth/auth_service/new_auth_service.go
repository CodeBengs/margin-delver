package authservice

import (
	"context"

	moduleLib "margin-delver/lib"
	authdto "margin-delver/modules/auth/auth_dto"
	authrepository "margin-delver/modules/auth/auth_entity/auth_repository"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, request *authdto.LoginRequest) (*authdto.LoginResponse, error)
}

type AuthService struct {
	authRepository authrepository.AuthRepositoryInterface
	log            *moduleLib.BaseLog
}

func NewService(
	authRepository authrepository.AuthRepositoryInterface,
	log *moduleLib.BaseLog,
) AuthServiceInterface {
	return &AuthService{
		authRepository: authRepository,
		log:            log,
	}
}
