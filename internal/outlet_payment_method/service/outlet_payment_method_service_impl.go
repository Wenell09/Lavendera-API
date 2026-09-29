package service

import (
	"context"
	"errors"
	"math"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/outlet/repository"
	"github.com/Wenell09/lavendera-api/internal/outlet_payment_method/dto"
	paymentMethodRepo "github.com/Wenell09/lavendera-api/internal/outlet_payment_method/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OutletPaymentMethodServiceImpl struct {
	PaymentMethodRepo paymentMethodRepo.OutletPaymentMethodRepository
	OutletRepo        repository.OutletRepository
	Logger            *logrus.Logger
}

func NewOutletPaymentMethodService(paymentMethodRepo paymentMethodRepo.OutletPaymentMethodRepository, outletRepo repository.OutletRepository, logger *logrus.Logger) OutletPaymentMethodService {
	return &OutletPaymentMethodServiceImpl{
		PaymentMethodRepo: paymentMethodRepo,
		OutletRepo:        outletRepo,
		Logger:            logger,
	}
}

func (s *OutletPaymentMethodServiceImpl) Create(ctx context.Context, req dto.OutletPaymentMethodRequest) (*dto.OutletPaymentMethodResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("outlet_id", req.OutletID)

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	outlet, err := s.OutletRepo.FindByID(ctx, req.OutletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet")
		return nil, err
	}

	exists2, err := s.PaymentMethodRepo.Exists(ctx, req.OutletID, req.Type, req.ProviderName, req.AccountNumber)
	if err != nil {
		logger.WithError(err).Error("failed to check existing outlet payment method")
		return nil, err
	}
	if exists2 {
		return nil, apperror.ConflictError{Msg: "outlet payment method already exists"}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	payment := models.OutletPaymentMethod{
		TenantID:      tenantID,
		OutletID:      req.OutletID,
		Type:          req.Type,
		ProviderName:  req.ProviderName,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		QRImageURL:    req.QRImageURL,
		IsActive:      isActive,
		Outlet:        *outlet,
	}

	if err := s.PaymentMethodRepo.Create(ctx, &payment); err != nil {
		logger.WithError(err).Error("failed to create outlet payment method")
		return nil, err
	}

	res := toOutletPaymentMethodResponse(&payment)
	return &res, nil
}

func (s *OutletPaymentMethodServiceImpl) FindAll(ctx context.Context, outletID uuid.UUID, filter dto.OutletPaymentMethodFilter) (*dto.OutletPaymentMethodListResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("outlet_id", outletID)
	filter.SetDefault()

	payments, total, err := s.PaymentMethodRepo.FindAll(ctx, outletID, filter)
	if err != nil {
		logger.WithError(err).Error("failed to find outlet payment methods")
		return nil, err
	}

	data := make([]dto.OutletPaymentMethodResponse, len(payments))
	for i, p := range payments {
		data[i] = toOutletPaymentMethodResponse(&p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.OutletPaymentMethodListResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *OutletPaymentMethodServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.OutletPaymentMethodResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("payment_id", id)

	payment, err := s.PaymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet payment method not found"}
		}
		logger.WithError(err).Error("failed to find outlet payment method")
		return nil, err
	}

	res := toOutletPaymentMethodResponse(payment)
	return &res, nil
}

func (s *OutletPaymentMethodServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateOutletPaymentMethodRequest) (*dto.OutletPaymentMethodResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("payment_id", id)

	payment, err := s.PaymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet payment method not found"}
		}
		logger.WithError(err).Error("failed to find outlet payment method")
		return nil, err
	}

	if req.Type != nil {
		payment.Type = *req.Type
	}
	if req.ProviderName != nil {
		payment.ProviderName = *req.ProviderName
	}
	if req.AccountNumber != nil {
		payment.AccountNumber = req.AccountNumber
	}
	if req.AccountName != nil {
		payment.AccountName = *req.AccountName
	}
	if req.QRImageURL != nil {
		payment.QRImageURL = req.QRImageURL
	}
	if req.IsActive != nil {
		payment.IsActive = *req.IsActive
	}

	if err := s.PaymentMethodRepo.Update(ctx, payment); err != nil {
		logger.WithError(err).Error("failed to update outlet payment method")
		return nil, err
	}

	res := toOutletPaymentMethodResponse(payment)
	return &res, nil
}

func (s *OutletPaymentMethodServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("payment_id", id)

	_, err := s.PaymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFoundError{Msg: "outlet payment method not found"}
		}
		logger.WithError(err).Error("failed to find outlet payment method")
		return err
	}

	if err := s.PaymentMethodRepo.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete outlet payment method")
		return err
	}

	return nil
}

func toOutletPaymentMethodResponse(m *models.OutletPaymentMethod) dto.OutletPaymentMethodResponse {
	res := dto.OutletPaymentMethodResponse{
		ID:            m.ID,
		OutletID:      m.OutletID,
		Outlet:        dto.OutletResponse{Name: m.Outlet.Name, Slug: m.Outlet.Slug},
		Type:          m.Type,
		ProviderName:  m.ProviderName,
		AccountNumber: m.AccountNumber,
		AccountName:   m.AccountName,
		QRImageURL:    m.QRImageURL,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	return res
}
