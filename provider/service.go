package provider

import (
	authprovider "margin-delver/modules/auth/auth_provider"

	"go.uber.org/fx"
)

func NewService() fx.Option {
	return fx.Provide(
		authprovider.NewServiceProvider,
	)
}
