package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	customerRepo "github.com/Wenell09/lavendera-api/internal/customer/repository"
	discountRepo "github.com/Wenell09/lavendera-api/internal/discount/repository"
	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/order/dto"
	"github.com/Wenell09/lavendera-api/internal/order/repository"
	outletRepo "github.com/Wenell09/lavendera-api/internal/outlet/repository"
	serviceRepo "github.com/Wenell09/lavendera-api/internal/service/repository"
	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/Wenell09/lavendera-api/internal/shared/apperror"
	"github.com/Wenell09/lavendera-api/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	Repository         repository.OrderRepository
	OutletRepository   outletRepo.OutletRepository
	CustomerRepository customerRepo.CustomerRepository
	ServiceRepository  serviceRepo.ServiceRepository
	DiscountRepository discountRepo.DiscountRepository
	Logger             *logrus.Logger
}

func NewOrderService(
	repository repository.OrderRepository,
	outletRepository outletRepo.OutletRepository,
	customerRepository customerRepo.CustomerRepository,
	serviceRepository serviceRepo.ServiceRepository,
	discountRepository discountRepo.DiscountRepository,
	logger *logrus.Logger,
) OrderService {
	return &OrderServiceImpl{
		Repository:         repository,
		OutletRepository:   outletRepository,
		CustomerRepository: customerRepository,
		ServiceRepository:  serviceRepository,
		DiscountRepository: discountRepository,
		Logger:             logger,
	}
}

// CreateCashierOrder implements [OrderService].
func (s *OrderServiceImpl) CreateCashierOrder(ctx context.Context, req dto.CreateCashierOrderRequest) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"outlet_id": req.OutletID,
	})

	tenantID, exists := appcontext.TenantIDFromContext(ctx)
	if !exists || tenantID == uuid.Nil {
		return nil, apperror.UnauthorizedError{Msg: "tenant_id not found"}
	}

	outlet, err := s.OutletRepository.FindByID(ctx, req.OutletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet")
		return nil, err
	}
	if !outlet.IsActive {
		return nil, apperror.ValidationError{Msg: "outlet is inactive"}
	}

	if req.CustomerID == uuid.Nil {
		return nil, apperror.ValidationError{Msg: "customer_id is required"}
	}
	if _, err := s.CustomerRepository.FindByID(ctx, req.CustomerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "customer not found"}
		}
		logger.WithError(err).Error("failed to find customer")
		return nil, err
	}

	//Process Items and calculate pricing for kiloan and satuan
	orderItems, subtotal, maxDurationDays, err := s.processOrderItems(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	//Calculate Discount
	discountAmount, err := s.calculateDiscount(ctx, req.DiscountID, subtotal)
	if err != nil {
		return nil, err
	}

	total := subtotal - discountAmount
	if total < 0 {
		total = 0
	}

	// 5. Generate unique tracking token and order number
	trackingToken, err := s.generateUniqueTrackingToken(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to generate tracking token")
		return nil, err
	}

	orderNumber, err := s.generateUniqueOrderNumber(ctx, tenantID)
	if err != nil {
		logger.WithError(err).Error("failed to generate order number")
		return nil, err
	}

	// 6. Creator ID (from JWT)
	var createdBy *uuid.UUID
	if uid, ok := appcontext.UserIDFromContext(ctx); ok && uid != uuid.Nil {
		exists, _ := s.Repository.UserExists(ctx, uid)
		if exists {
			createdBy = &uid
		}
	}

	// 7. Status
	status := "DIPROSES"
	if req.Status != "" {
		status = req.Status
	}

	orderID := uuid.New()
	for i := range orderItems {
		orderItems[i].OrderID = orderID
	}

	var paidAmount int64 = 0
	var payment *models.Payment

	if req.PaymentMethod != "" && strings.TrimSpace(req.PaymentMethod) != "" {
		method := strings.ToLower(strings.TrimSpace(req.PaymentMethod))
		refNum, err := GenerateReferenceNumber()
		if err != nil {
			logger.WithError(err).Error("failed to generate payment reference number")
			return nil, err
		}

		var paymentAmount int64
		if method == "cash" {
			if *req.PaidAmount < total {
				return nil, apperror.ValidationError{Msg: "paid_amount must be greater than or equal to total price for cash payment"}
			}
			paymentAmount = *req.PaidAmount
		} else {
			paymentAmount = total
		}
		paidAmount = paymentAmount

		now := time.Now()
		payment = &models.Payment{
			ID:              uuid.New(),
			OutletID:        outlet.ID,
			OrderID:         orderID,
			ReceivedBy:      createdBy,
			Amount:          paymentAmount,
			PaymentMethod:   method,
			Status:          "verified",
			ReferenceNumber: refNum,
			PaidAt:          &now,
			VerifiedAt:      &now,
		}
	}

	order := &models.Order{
		ID:             orderID,
		TenantID:       tenantID,
		OutletID:       outlet.ID,
		CustomerID:     req.CustomerID,
		DiscountID:     req.DiscountID,
		CreatedBy:      createdBy,
		OrderNumber:    orderNumber,
		TrackingToken:  trackingToken,
		Status:         status,
		Subtotal:       subtotal,
		DiscountAmount: discountAmount,
		Total:          total,
		PaidAmount:     paidAmount,
		Notes:          req.Notes,
		EstimatedDone:  time.Now().Add(time.Duration(maxDurationDays) * 24 * time.Hour),
		Items:          orderItems,
	}

	if err := s.Repository.CreateWithPayment(ctx, order, payment); err != nil {
		logger.WithError(err).Error("failed to create order with payment")
		return nil, err
	}

	// Reload created order with preloaded relations
	createdOrder, err := s.Repository.FindByID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to reload created order")
		return s.mapOrderToResponse(order), nil
	}

	logger.WithFields(logrus.Fields{
		"order_id":       createdOrder.ID,
		"order_number":   createdOrder.OrderNumber,
		"tracking_token": createdOrder.TrackingToken,
		"has_payment":    payment != nil,
	}).Info("cashier order created successfully")

	return s.mapOrderToResponse(createdOrder), nil
}

// CreatePublicOrder implements [OrderService].
func (s *OrderServiceImpl) CreatePublicOrder(ctx context.Context, req dto.CreatePublicOrderRequest) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"outlet_slug": req.OutletSlug,
		"phone":       req.Phone,
	})

	// 1. Find Outlet by slug (public finder without tenant context)
	outlet, err := s.OutletRepository.FindBySlugPublic(ctx, req.OutletSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "outlet not found"}
		}
		logger.WithError(err).Error("failed to find outlet by slug")
		return nil, err
	}
	if !outlet.IsActive {
		return nil, apperror.ValidationError{Msg: "outlet is inactive"}
	}
	if !outlet.IsPublicOrderEnabled {
		return nil, apperror.ValidationError{Msg: "public order is not enabled for this outlet"}
	}

	tenantID := outlet.TenantID
	ctx = appcontext.WithTenantID(ctx, tenantID)

	// 2. Resolve Customer by name AND phone or create new one
	phone := strings.TrimSpace(req.Phone)
	name := strings.TrimSpace(req.Name)

	var customerID uuid.UUID
	existingCustomer, err := s.CustomerRepository.FindByNameAndPhone(ctx, tenantID, name, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCustomer := &models.Customer{
				TenantID: tenantID,
				Name:     name,
				Phone:    phone,
				Address:  req.Address,
			}
			if err := s.CustomerRepository.Create(ctx, newCustomer); err != nil {
				logger.WithError(err).Error("failed to create customer for public order")
				return nil, err
			}
			customerID = newCustomer.ID
		} else {
			logger.WithError(err).Error("failed to check existing customer by name and phone")
			return nil, err
		}
	} else {
		customerID = existingCustomer.ID
	}

	// 3. Process items and calculate pricing for kiloan and satuan
	orderItems, subtotal, maxDurationDays, err := s.processOrderItems(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	// 4. Calculate discount
	discountAmount, err := s.calculateDiscount(ctx, req.DiscountID, subtotal)
	if err != nil {
		return nil, err
	}

	total := subtotal - discountAmount
	if total < 0 {
		total = 0
	}

	// 5. Generate tracking token and order number
	trackingToken, err := s.generateUniqueTrackingToken(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to generate tracking token")
		return nil, err
	}

	orderNumber, err := s.generateUniqueOrderNumber(ctx, tenantID)
	if err != nil {
		logger.WithError(err).Error("failed to generate order number")
		return nil, err
	}

	// 6. Estimated Done
	est := time.Now().Add(time.Duration(maxDurationDays) * 24 * time.Hour)

	orderID := uuid.New()
	for i := range orderItems {
		orderItems[i].OrderID = orderID
	}

	var paidAmount int64 = 0
	var payment *models.Payment

	if req.PaymentMethod != "" && strings.TrimSpace(req.PaymentMethod) != "" {
		method := strings.ToLower(strings.TrimSpace(req.PaymentMethod))
		refNum, err := GenerateReferenceNumber()
		if err != nil {
			logger.WithError(err).Error("failed to generate payment reference number")
			return nil, err
		}
		var paymentAmount int64
		if method == "cash" {
			if req.PaidAmount != nil && *req.PaidAmount != 0 && *req.PaidAmount < total {
				return nil, apperror.ValidationError{Msg: "paid_amount must be greater than or equal to total price for cash payment"}
			}
			if req.PaidAmount != nil {
				paymentAmount = *req.PaidAmount
			} else {
				paymentAmount = total
			}
		} else {
			paymentAmount = total
		}
		// keep order paid_amount zero until verification

		payment = &models.Payment{
			ID:              uuid.New(),
			OutletID:        outlet.ID,
			OrderID:         orderID,
			ReceivedBy:      nil,
			Amount:          paymentAmount,
			PaymentMethod:   method,
			Status:          "pending",
			ReferenceNumber: refNum,
		}
	}

	order := &models.Order{
		ID:             orderID,
		TenantID:       tenantID,
		OutletID:       outlet.ID,
		CustomerID:     customerID,
		DiscountID:     req.DiscountID,
		CreatedBy:      nil,
		OrderNumber:    orderNumber,
		TrackingToken:  trackingToken,
		Status:         "MENUNGGU_KONFIRMASI",
		Subtotal:       subtotal,
		DiscountAmount: discountAmount,
		Total:          total,
		PaidAmount:     paidAmount,
		Notes:          req.Notes,
		EstimatedDone:  est,
		Items:          orderItems,
	}

	if err := s.Repository.CreateWithPayment(ctx, order, payment); err != nil {
		logger.WithError(err).Error("failed to create public order with payment")
		return nil, err
	}

	// Reload created order with preloaded relations
	createdOrder, err := s.Repository.FindByID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to reload created public order")
		return s.mapOrderToResponse(order), nil
	}

	logger.WithFields(logrus.Fields{
		"order_id":       createdOrder.ID,
		"order_number":   createdOrder.OrderNumber,
		"tracking_token": createdOrder.TrackingToken,
		"has_payment":    payment != nil,
	}).Info("public order created successfully")

	return s.mapOrderToResponse(createdOrder), nil
}

// FindAll implements [OrderService].
func (s *OrderServiceImpl) FindAll(ctx context.Context, filter dto.OrderFilter) (*dto.OrderListResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx)

	filter.SetDefault()

	orders, total, err := s.Repository.FindAll(ctx, filter)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve orders")
		return nil, err
	}

	orderResponses := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		orderResponses[i] = *s.mapOrderToResponse(&orders[i])
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.OrderListResponse{
		Data: orderResponses,
		Pagination: dto.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// FindByID implements [OrderService].
func (s *OrderServiceImpl) FindByID(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("id", id)

	order, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "order not found"}
		}
		logger.WithError(err).Error("failed to find order by id")
		return nil, err
	}

	return s.mapOrderToResponse(order), nil
}

// FindByTrackingToken implements [OrderService].
func (s *OrderServiceImpl) FindByTrackingToken(ctx context.Context, trackingToken string) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithField("tracking_token", trackingToken)

	token := strings.TrimSpace(trackingToken)
	order, err := s.Repository.FindByTrackingToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "order with tracking token not found"}
		}
		logger.WithError(err).Error("failed to find order by tracking token")
		return nil, err
	}

	return s.mapOrderToResponse(order), nil
}

// UpdateStatus implements [OrderService].
func (s *OrderServiceImpl) UpdateStatus(ctx context.Context, id uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"id":     id,
		"status": req.Status,
	})

	order, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "order not found"}
		}
		logger.WithError(err).Error("failed to find order for status update")
		return nil, err
	}

	if order.Status == "DIBATALKAN" {
		return nil, apperror.ValidationError{Msg: "cannot update status of canceled order"}
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == "SELESAI" {
		updates["completed_at"] = &now
	}

	if err := s.Repository.Update(ctx, order.ID, updates); err != nil {
		logger.WithError(err).Error("failed to update order status")
		return nil, err
	}

	updatedOrder, err := s.Repository.FindByID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to reload order after status update")
		return nil, err
	}

	logger.Info("order status updated successfully")
	return s.mapOrderToResponse(updatedOrder), nil
}

// CancelOrder implements [OrderService].
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"id": id,
	})

	order, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "order not found"}
		}
		logger.WithError(err).Error("failed to find order for cancellation")
		return nil, err
	}

	if order.Status == "DIBATALKAN" {
		return nil, apperror.ValidationError{Msg: "order is already canceled"}
	}
	if order.Status == "SELESAI" || order.Status == "DIAMBIL" {
		return nil, apperror.ValidationError{Msg: "cannot cancel completed or picked up order"}
	}

	now := time.Now()
	orderUpdates := map[string]interface{}{
		"status":      "DIBATALKAN",
		"canceled_at": &now,
		"paid_amount": 0,
	}

	paymentUpdates := map[string]interface{}{
		"status":      "rejected",
		"verified_at": nil,
	}

	if err := s.Repository.CancelOrderWithPayments(ctx, order.ID, orderUpdates, paymentUpdates); err != nil {
		logger.WithError(err).Error("failed to cancel order and update payments")
		return nil, err
	}

	updatedOrder, err := s.Repository.FindByID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to reload order after cancellation")
		return nil, err
	}

	logger.Info("order canceled successfully")
	return s.mapOrderToResponse(updatedOrder), nil
}

// ConfirmPayment implements [OrderService].
func (s *OrderServiceImpl) ConfirmPayment(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"order_id":   orderID,
		"payment_id": paymentID,
	})

	// 1. Verify order exists
	order, err := s.Repository.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "order not found"}
		}
		logger.WithError(err).Error("failed to find order for payment confirm")
		return nil, err
	}

	// 2. Verify payment exists and belongs to order
	payment, err := s.Repository.FindPaymentByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "payment not found"}
		}
		logger.WithError(err).Error("failed to find payment")
		return nil, err
	}

	if payment.OrderID != order.ID {
		return nil, apperror.NotFoundError{Msg: "payment does not belong to this order"}
	}

	now := time.Now()
	paymentUpdates := map[string]interface{}{
		"status":      "verified",
		"verified_at": &now,
		"paid_at":     &now,
	}
	if uid, ok := appcontext.UserIDFromContext(ctx); ok && uid != uuid.Nil {
		paymentUpdates["received_by"] = &uid
	}

	currentVerifiedSum, err := s.Repository.SumVerifiedPaymentsByOrderID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to calculate sum of verified payments")
		return nil, err
	}

	var newPaidAmount int64
	if payment.Status == "verified" {
		newPaidAmount = currentVerifiedSum
	} else {
		newPaidAmount = currentVerifiedSum + payment.Amount
	}

	if err := s.Repository.UpdatePaymentAndOrderPaidAmount(ctx, payment.ID, paymentUpdates, order.ID, newPaidAmount); err != nil {
		logger.WithError(err).Error("failed to confirm payment")
		return nil, err
	}

	updatedOrder, err := s.Repository.FindByID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to reload order after payment confirm")
		return nil, err
	}

	logger.Info("payment confirmed successfully")
	return s.mapOrderToResponse(updatedOrder), nil
}

// InputCashPayment implements [OrderService].
func (s *OrderServiceImpl) InputCashPayment(ctx context.Context, id uuid.UUID, req dto.UpdateOrderPaidAmountRequest) (*dto.OrderResponse, error) {
	logger := utils.LogWithContext(s.Logger, ctx).WithFields(logrus.Fields{
		"order_id":    id,
		"paid_amount": req.PaidAmount,
	})

	order, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError{Msg: "order not found"}
		}
		return nil, err
	}

	if order.Status == "DIBATALKAN" {
		return nil, apperror.ValidationError{Msg: "cannot process payment for canceled order"}
	}

	if order.Payment == nil || order.Payment.Status != "pending" || order.Payment.PaymentMethod != "cash" {
		return nil, apperror.NotFoundError{Msg: "no pending cash payment found for this order"}
	}
	pendingPayment := order.Payment

	if req.PaidAmount < order.Total {
		return nil, apperror.ValidationError{Msg: "paid_amount must be greater than or equal to total price"}
	}

	now := time.Now()
	paymentUpdates := map[string]interface{}{
		"status":      "verified",
		"verified_at": &now,
		"paid_at":     &now,
		"amount":      req.PaidAmount,
	}
	if uid, ok := appcontext.UserIDFromContext(ctx); ok && uid != uuid.Nil {
		paymentUpdates["received_by"] = &uid
	}

	if err := s.Repository.UpdatePaymentAndOrderPaidAmount(ctx, pendingPayment.ID, paymentUpdates, order.ID, req.PaidAmount); err != nil {
		logger.WithError(err).Error("failed to input cash payment")
		return nil, err
	}

	// Also update status to DIPROSES if it was MENUNGGU_KONFIRMASI
	if order.Status == "MENUNGGU_KONFIRMASI" {
		if err := s.Repository.Update(ctx, order.ID, map[string]interface{}{"status": "DIPROSES"}); err != nil {
			logger.WithError(err).Error("failed to update order status to DIPROSES")
			// continue as payment is already recorded
		}
	}

	updatedOrder, err := s.Repository.FindByID(ctx, order.ID)
	if err != nil {
		logger.WithError(err).Error("failed to reload order after cash payment")
		return nil, err
	}

	return s.mapOrderToResponse(updatedOrder), nil
}

// processOrderItems processes items, verifies services directly from system,
// validates uniqueness of services in the item list, and calculates pricing:
// - kiloan (unit "gram"): subtotal = round(unitPrice * (quantity / 1000.0))
// - satuan (non-gram such as "pcs", "pair"): subtotal = unitPrice * quantity
func (s *OrderServiceImpl) processOrderItems(ctx context.Context, items []dto.CreateOrderItemRequest) ([]models.OrderItem, int64, int64, error) {
	if len(items) == 0 {
		return nil, 0, 0, apperror.ValidationError{Msg: "order must contain at least one item"}
	}

	seenServices := make(map[uuid.UUID]bool)
	var orderItems []models.OrderItem
	var subtotal int64
	var maxDurationDays int64 = 1

	for _, item := range items {
		if seenServices[item.ServiceID] {
			return nil, 0, 0, apperror.ValidationError{
				Msg: fmt.Sprintf("duplicate service with id %s in items", item.ServiceID),
			}
		}
		seenServices[item.ServiceID] = true

		service, err := s.ServiceRepository.FindByID(ctx, item.ServiceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, 0, apperror.NotFoundError{
					Msg: fmt.Sprintf("service with id %s not found", item.ServiceID),
				}
			}
			return nil, 0, 0, err
		}

		if !service.IsActive {
			return nil, 0, 0, apperror.ValidationError{
				Msg: fmt.Sprintf("service '%s' is inactive", service.Name),
			}
		}

		checkQty := item.Quantity
		if strings.ToLower(service.Unit) == "gram" {
			checkQty = item.Quantity * 1000
		}

		if service.MinQuantity > 0 && checkQty < float64(service.MinQuantity) {
			return nil, 0, 0, apperror.ValidationError{
				Msg: fmt.Sprintf("quantity for service '%s' must be at least %d %s", service.Name, service.MinQuantity, service.Unit),
			}
		}

		unitPrice := service.Price
		var itemSubtotal int64
		var storedQuantity int64

		if strings.ToLower(service.Unit) == "gram" {
			itemSubtotal = int64(math.Round(float64(unitPrice) * item.Quantity))
			storedQuantity = int64(math.Round(item.Quantity * 1000))
		} else {
			itemSubtotal = int64(math.Round(float64(unitPrice) * item.Quantity))
			storedQuantity = int64(math.Round(item.Quantity))
		}

		subtotal += itemSubtotal

		if service.DurationDays > maxDurationDays {
			maxDurationDays = service.DurationDays
		}

		orderItems = append(orderItems, models.OrderItem{
			ID:        uuid.New(),
			ServiceID: service.ID,
			Price:     unitPrice,
			Quantity:  storedQuantity,
			Subtotal:  itemSubtotal,
		})
	}

	return orderItems, subtotal, maxDurationDays, nil
}

// calculateDiscount calculates discount amount based on discount type and subtotal.
func (s *OrderServiceImpl) calculateDiscount(ctx context.Context, discountID *uuid.UUID, subtotal int64) (int64, error) {
	if discountID == nil || *discountID == uuid.Nil {
		return 0, nil
	}

	discount, err := s.DiscountRepository.FindByID(ctx, *discountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, apperror.NotFoundError{Msg: "discount not found"}
		}
		return 0, err
	}

	if !discount.IsActive {
		return 0, apperror.ValidationError{Msg: fmt.Sprintf("discount '%s' is inactive", discount.Name)}
	}

	var discountAmount int64
	if discount.Type == "percentage" {
		discountAmount = (subtotal * discount.Value) / 100
	} else {
		discountAmount = discount.Value
	}

	if discountAmount > subtotal {
		discountAmount = subtotal
	}

	return discountAmount, nil
}

// generateUniqueTrackingToken generates a simple, collision-free tracking token.
func (s *OrderServiceImpl) generateUniqueTrackingToken(ctx context.Context) (string, error) {
	for i := 0; i < 10; i++ {
		token, err := GenerateTrackingToken()
		if err != nil {
			return "", err
		}
		exists, err := s.Repository.ExistsByTrackingToken(ctx, token)
		if err != nil {
			return "", err
		}
		if !exists {
			return token, nil
		}
	}
	return "", errors.New("unable to generate unique tracking token")
}

// generateUniqueOrderNumber generates a unique order number per tenant.
func (s *OrderServiceImpl) generateUniqueOrderNumber(ctx context.Context, tenantID uuid.UUID) (string, error) {
	for i := 0; i < 10; i++ {
		orderNumber, err := GenerateOrderNumber()
		if err != nil {
			return "", err
		}
		exists, err := s.Repository.ExistsByOrderNumber(ctx, tenantID, orderNumber)
		if err != nil {
			return "", err
		}
		if !exists {
			return orderNumber, nil
		}
	}
	return "", errors.New("unable to generate unique order number")
}

// mapOrderToResponse converts models.Order to dto.OrderResponse.
func (s *OrderServiceImpl) mapOrderToResponse(order *models.Order) *dto.OrderResponse {
	if order == nil {
		return nil
	}

	paymentStatus := "UNPAID"
	if order.PaidAmount >= totalNonNegative(order.Total) && totalNonNegative(order.Total) > 0 {
		paymentStatus = "PAID"
	} else if order.PaidAmount > 0 {
		paymentStatus = "PARTIAL"
	}

	var changeAmount int64 = 0
	if order.PaidAmount > totalNonNegative(order.Total) {
		changeAmount = order.PaidAmount - totalNonNegative(order.Total)
	}

	resp := &dto.OrderResponse{
		ID:             order.ID,
		TenantID:       order.TenantID,
		OrderNumber:    order.OrderNumber,
		TrackingToken:  order.TrackingToken,
		Status:         order.Status,
		PaymentStatus:  paymentStatus,
		Subtotal:       order.Subtotal,
		DiscountAmount: order.DiscountAmount,
		TotalAmount:    totalNonNegative(order.Total),
		PaidAmount:     order.PaidAmount,
		ChangeAmount:   changeAmount,
		Notes:          order.Notes,
		EstimatedDone:  order.EstimatedDone,
		CompletedAt:    order.CompletedAt,
		CanceledAt:     order.CanceledAt,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

	if order.Outlet.ID != uuid.Nil {
		resp.Outlet = dto.OutletSummary{
			ID:      order.Outlet.ID,
			Name:    order.Outlet.Name,
			Phone:   order.Outlet.Phone,
			Address: order.Outlet.Address,
		}
	}

	if order.Customer.ID != uuid.Nil {
		resp.Customer = dto.CustomerSummary{
			ID:      order.Customer.ID,
			Name:    order.Customer.Name,
			Phone:   order.Customer.Phone,
			Address: order.Customer.Address,
		}
	}

	if order.Discount != nil && order.Discount.ID != uuid.Nil {
		resp.Discount = &dto.DiscountSummary{
			ID:    order.Discount.ID,
			Name:  order.Discount.Name,
			Type:  order.Discount.Type,
			Value: order.Discount.Value,
		}
	}

	if order.Creator != nil && order.Creator.ID != uuid.Nil {
		resp.Creator = &dto.UserSummary{
			ID:   order.Creator.ID,
			Name: order.Creator.Name,
		}
	}

	if len(order.Items) > 0 {
		items := make([]dto.OrderItemResponse, len(order.Items))
		for i, item := range order.Items {
			qty := float64(item.Quantity)
			unit := "pcs"
			if item.Service.ID != uuid.Nil {
				unit = item.Service.Unit
			}

			if strings.ToLower(unit) == "gram" {
				qty = float64(item.Quantity) / 1000.0
				unit = "kg"
			}

			itemResp := dto.OrderItemResponse{
				ID:       item.ID,
				Price:    item.Price,
				Quantity: qty,
				Unit:     unit,
				Subtotal: item.Subtotal,
			}
			if item.Service.ID != uuid.Nil {
				srv := dto.ServiceSummary{
					ID:           item.Service.ID,
					Name:         item.Service.Name,
					Price:        item.Service.Price,
					DurationDays: item.Service.DurationDays,
				}
				if item.Service.Category.ID != uuid.Nil {
					srv.Category = dto.ServiceCategorySummary{
						ID:   item.Service.Category.ID,
						Name: item.Service.Category.Name,
					}
				}
				itemResp.Service = srv
			}
			items[i] = itemResp
		}
		resp.Items = items
	}

	if order.Payment != nil && order.Payment.ID != uuid.Nil {
		var refNum *string
		if order.Payment.ReferenceNumber != "" {
			refNum = &order.Payment.ReferenceNumber
		}
		resp.Payment = &dto.PaymentResponse{
			ID:              order.Payment.ID,
			Amount:          order.Payment.Amount,
			PaymentMethod:   order.Payment.PaymentMethod,
			Status:          order.Payment.Status,
			ReferenceNumber: refNum,
			Notes:           order.Payment.Notes,
			PaidAt:          order.Payment.PaidAt,
			VerifiedAt:      order.Payment.VerifiedAt,
			CreatedAt:       order.Payment.CreatedAt,
			UpdatedAt:       order.Payment.UpdatedAt,
		}
	}

	return resp
}

func totalNonNegative(total int64) int64 {
	if total < 0 {
		return 0
	}
	return total
}
