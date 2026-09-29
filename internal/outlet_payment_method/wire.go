package outlet_payment_method

import (
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/controller"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/repository"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewOutletPaymentMethodRepository,
	service.NewOutletPaymentMethodService,
	controller.NewOutletPaymentMethodController,
)
