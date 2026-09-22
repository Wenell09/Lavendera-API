package discount

import (
	"github.com/Wenell09/lavendera-api/internal/discount/controller"
	"github.com/Wenell09/lavendera-api/internal/discount/repository"
	"github.com/Wenell09/lavendera-api/internal/discount/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewDiscountRepository,
	service.NewDiscountService,
	controller.NewDiscountController,
)
