package main

import (
	"margin-delver/cmd"
	"margin-delver/lib"
	"margin-delver/provider"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			lib.NewAppConfig,
			lib.NewBaseLog,
			lib.NewDatabase,
		),
		provider.NewRepository(),
		provider.NewService(),
		provider.NewHandler(),
		provider.NewSeeder(),
		fx.Invoke(
			lib.RunMigrations,
			cmd.NewServer,
		),
	).Run()
}
