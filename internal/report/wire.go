package report

import (
	"github.com/Wenell09/lavendera-api/internal/report/controller"
	"github.com/Wenell09/lavendera-api/internal/report/repository"
	"github.com/Wenell09/lavendera-api/internal/report/service"
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	repository.NewReportRepository,
	service.NewReportService,
	controller.NewReportController,
)
