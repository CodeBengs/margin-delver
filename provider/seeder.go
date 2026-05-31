package provider

import (
	"context"

	moduleLib "margin-delver/lib"
	authconstant "margin-delver/modules/auth/auth_constant"
	authmodel "margin-delver/modules/auth/auth_entity/auth_model"
	authrepository "margin-delver/modules/auth/auth_entity/auth_repository"

	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

func NewSeeder() fx.Option {
	return fx.Invoke(SeedDefaultUser)
}

func SeedDefaultUser(
	cfg *moduleLib.AppConfig,
	log *moduleLib.BaseLog,
	authRepository authrepository.AuthRepositoryInterface,
) error {
	if !cfg.DBSeedDefaultUser {
		log.SugarLog().Info("default auth user seed skipped")
		return nil
	}

	if cfg.AuthDefaultUsername == "" || cfg.AuthDefaultPassword == "" {
		log.SugarLog().Info("default auth user seed skipped")
		return nil
	}

	total, err := authRepository.Count(context.Background())
	if err != nil {
		log.SugarLog().Errorf("%s: %v", authconstant.ErrCodeDefaultUserSeedFail, err)
		return err
	}

	if total > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AuthDefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.SugarLog().Errorf("%s: %v", authconstant.ErrCodeDefaultUserSeedFail, err)
		return err
	}

	return authRepository.Create(context.Background(), &authmodel.User{
		Username:     cfg.AuthDefaultUsername,
		PasswordHash: string(passwordHash),
		Name:         cfg.AuthDefaultName,
		FlagActive:   true,
	})
}
