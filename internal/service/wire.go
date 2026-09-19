package service

import (
	"github.com/Wenell09/lavendera-api/internal/service/controller"
	"github.com/Wenell09/lavendera-api/internal/service/repository"
	"github.com/Wenell09/lavendera-api/internal/service/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewServiceRepository,
	service.NewService,
	controller.NewServiceController,
)
