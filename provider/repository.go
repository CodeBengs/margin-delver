package provider

import (
	authprovider "margin-delver/modules/auth/auth_provider"

	"go.uber.org/fx"
)

func NewRepository() fx.Option {
	return fx.Provide(
		authprovider.NewRepositoryProvider,
	)
}
