package service_category

import (
	"github.com/Wenell09/lavendera-api/internal/service_category/controller"
	"github.com/Wenell09/lavendera-api/internal/service_category/repository"
	"github.com/Wenell09/lavendera-api/internal/service_category/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewServiceCategoryRepository,
	service.NewServiceCategoryService,
	controller.NewServiceCategoryController,
)
