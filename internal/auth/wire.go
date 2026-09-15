package auth

import (
	"github.com/Wenell09/lavendera-api/internal/auth/controller"
	"github.com/Wenell09/lavendera-api/internal/auth/repository"
	"github.com/Wenell09/lavendera-api/internal/auth/service"
	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewAuthRepository,
	service.NewAuthService,
	controller.NewAuthController,

	config.NewValidator,
	config.NewLogger,
	config.NewJWTConfig,
)
