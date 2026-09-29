package customer

import (
	"github.com/Wenell09/lavendera-api/internal/customer/controller"
	"github.com/Wenell09/lavendera-api/internal/customer/repository"
	"github.com/Wenell09/lavendera-api/internal/customer/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewCustomerRepository,
	service.NewCustomerService,
	controller.NewCustomerController,
)
