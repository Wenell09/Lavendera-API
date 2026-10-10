package order

import (
	"github.com/Wenell09/lavendera-api/internal/order/controller"
	"github.com/Wenell09/lavendera-api/internal/order/repository"
	"github.com/Wenell09/lavendera-api/internal/order/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewOrderRepository,
	service.NewOrderService,
	controller.NewOrderController,
)
