package authprovider

import (
	moduleLib "margin-delver/lib"
	authrepository "margin-delver/modules/auth/auth_entity/auth_repository"
	authhandler "margin-delver/modules/auth/auth_handler"
	authservice "margin-delver/modules/auth/auth_service"

	"gorm.io/gorm"
)

func NewRepositoryProvider(db *gorm.DB) authrepository.AuthRepositoryInterface {
	return authrepository.NewRepository(db)
}

func NewServiceProvider(
	authRepository authrepository.AuthRepositoryInterface,
	log *moduleLib.BaseLog,
) authservice.AuthServiceInterface {
	return authservice.NewService(authRepository, log)
}

func NewHandlerProvider(
	authService authservice.AuthServiceInterface,
	log *moduleLib.BaseLog,
) *authhandler.Handler {
	return authhandler.NewHandler(authService, log)
}

func InitializeAuthHandler(
	log *moduleLib.BaseLog,
	cfg *moduleLib.AppConfig,
	db *gorm.DB,
) *authhandler.Handler {
	authRepository := NewRepositoryProvider(db)
	authService := NewServiceProvider(authRepository, log)

	return NewHandlerProvider(authService, log)
}
