package outlet

import (
	"github.com/Wenell09/lavendera-api/internal/outlet/controller"
	"github.com/Wenell09/lavendera-api/internal/outlet/repository"
	"github.com/Wenell09/lavendera-api/internal/outlet/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewOutletRepository,
	service.NewOutletService,
	controller.NewOutletController,
)
