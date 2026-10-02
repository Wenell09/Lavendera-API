//go:build wireinject

package app

import (
	"github.com/Wenell09/lavendera-api/internal/auth"
	"github.com/Wenell09/lavendera-api/internal/customer"
	"github.com/Wenell09/lavendera-api/internal/discount"
	"github.com/Wenell09/lavendera-api/internal/outlet"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method"
	"github.com/Wenell09/lavendera-api/internal/outlet_user"
	"github.com/Wenell09/lavendera-api/internal/service"
	"github.com/Wenell09/lavendera-api/internal/tenant"
	"github.com/Wenell09/lavendera-api/internal/service_category"
	"github.com/Wenell09/lavendera-api/internal/shared/applogger"
	"github.com/Wenell09/lavendera-api/internal/shared/appvalidator"
	"github.com/Wenell09/lavendera-api/internal/user"
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		auth.WireSet,
		service_category.WireSet,
		service.WireSet,
		discount.WireSet,
		customer.WireSet,
		outlet.WireSet,
		outlet_payment_method.WireSet,
		user.WireSet,
		outlet_user.WireSet,
		tenant.WireSet,
		appvalidator.NewValidator,
		applogger.NewLogger,
		NewDB,
		NewRouter,
		NewApp,
	)

	return nil, nil
}
