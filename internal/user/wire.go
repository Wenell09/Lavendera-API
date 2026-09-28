package user

import (
	"github.com/Wenell09/lavendera-api/internal/user/controller"
	"github.com/Wenell09/lavendera-api/internal/user/repository"
	"github.com/Wenell09/lavendera-api/internal/user/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	controller.NewUserController,
)
