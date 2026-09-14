//go:build wireinject

package app

import (
	"github.com/Wenell09/lavendera-api/internal/auth/controller"
	"github.com/Wenell09/lavendera-api/internal/auth/repository"
	"github.com/Wenell09/lavendera-api/internal/auth/service"
	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/google/wire"
)

var AuthWireSet = wire.NewSet(
	repository.NewAuthRepository,
	service.NewAuthService,
	controller.NewAuthController,

	config.NewValidator,
	config.NewLogger,
	config.NewJWTConfig,
)

func InitializeApp() (*App, error) {
	wire.Build(
		AuthWireSet,

		NewDB,
		NewRouter,
		NewApp,
	)

	return nil, nil
}
