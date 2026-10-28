package profile

import (
	"github.com/Wenell09/lavendera-api/internal/profile/controller"
	"github.com/Wenell09/lavendera-api/internal/profile/repository"
	"github.com/Wenell09/lavendera-api/internal/profile/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewProfileRepository,
	service.NewProfileService,
	controller.NewProfileController,
)
