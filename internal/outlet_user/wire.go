package outlet_user

import (
	"github.com/Wenell09/lavendera-api/internal/outlet_user/controller"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/repository"
	"github.com/Wenell09/lavendera-api/internal/outlet_user/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewOutletUserRepository,
	service.NewOutletUserService,
	controller.NewOutletUserController,
)
