//go:build wireinject

package app

import (
	"github.com/Wenell09/lavendera-api/internal/auth"
	"github.com/Wenell09/lavendera-api/internal/service_category"
	"github.com/Wenell09/lavendera-api/internal/service"
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		auth.WireSet,
		service.WireSet,
		servicecategory.WireSet,
		NewDB,
		NewRouter,
		NewApp,
	)

	return nil, nil
}
