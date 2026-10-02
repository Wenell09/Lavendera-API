package tenant

import (
	"github.com/Wenell09/lavendera-api/internal/tenant/controller"
	"github.com/Wenell09/lavendera-api/internal/tenant/repository"
	"github.com/Wenell09/lavendera-api/internal/tenant/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewTenantRepository,
	service.NewTenantService,
	controller.NewTenantController,
)
