package repository

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/models"
	"github.com/Wenell09/lavendera-api/internal/order/dto"
	"github.com/Wenell09/lavendera-api/internal/shared/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{DB: db}
}

// CreateWithPayment implements [OrderRepository].
func (r *OrderRepositoryImpl) CreateWithPayment(ctx context.Context, order *models.Order, payment *models.Payment) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		if payment != nil {
			if err := tx.Create(payment).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindAll implements [OrderRepository].
func (r *OrderRepositoryImpl) FindAll(ctx context.Context, filter dto.OrderFilter) ([]models.Order, int64, error) {
	orders := []models.Order{}
	var total int64

	query := r.DB.WithContext(ctx).Scopes(database.TenantScope(ctx)).Model(&models.Order{}).
		Preload("Outlet").
		Preload("Customer").
		Preload("Discount").
		Preload("Creator").
		Preload("Items").
		Preload("Items.Service").
		Preload("Items.Service.Category").
		Preload("Payments")

	if filter.Search != "" {
		query = query.Where("order_number ILIKE ? OR tracking_token ILIKE ?", filter.Search+"%", filter.Search+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OutletID != nil {
		query = query.Where("outlet_id = ?", filter.OutletID)
	}
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// FindByID implements [OrderRepository].
func (r *OrderRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	if err := r.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Preload("Outlet").
		Preload("Customer").
		Preload("Discount").
		Preload("Creator").
		Preload("Items").
		Preload("Items.Service").
		Preload("Items.Service.Category").
		Preload("Payments").
		First(order, id).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// FindByTrackingToken implements [OrderRepository].
func (r *OrderRepositoryImpl) FindByTrackingToken(ctx context.Context, trackingToken string) (*models.Order, error) {
	order := &models.Order{}
	if err := r.DB.WithContext(ctx).
		Preload("Outlet").
		Preload("Customer").
		Preload("Discount").
		Preload("Creator").
		Preload("Items").
		Preload("Items.Service").
		Preload("Items.Service.Category").
		Preload("Payments").
		Where("tracking_token = ?", trackingToken).
		First(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// Update implements [OrderRepository].
func (r *OrderRepositoryImpl) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	if err := r.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Model(&models.Order{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// ExistsByTrackingToken implements [OrderRepository].
func (r *OrderRepositoryImpl) ExistsByTrackingToken(ctx context.Context, trackingToken string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("tracking_token = ?", trackingToken).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByOrderNumber implements [OrderRepository].
func (r *OrderRepositoryImpl) ExistsByOrderNumber(ctx context.Context, tenantID uuid.UUID, orderNumber string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Scopes(database.TenantScope(ctx)).
		Model(&models.Order{}).
		Where("order_number = ?", orderNumber).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindPaymentByID implements [OrderRepository].
func (r *OrderRepositoryImpl) FindPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	payment := &models.Payment{}
	if err := r.DB.WithContext(ctx).
		Scopes(database.TenantScopeByOrder(ctx)).
		First(payment, id).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

// SumVerifiedPaymentsByOrderID calculates sum of verified payments for an order.
func (r *OrderRepositoryImpl) SumVerifiedPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) (int64, error) {
	var totalPaid int64
	if err := r.DB.WithContext(ctx).Model(&models.Payment{}).
		Scopes(database.TenantScopeByOrder(ctx)).
		Where("order_id = ? AND status = ?", orderID, "verified").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid).Error; err != nil {
		return 0, err
	}
	return totalPaid, nil
}

// UpdatePaymentAndOrderPaidAmount updates payment and order paid_amount atomically.
func (r *OrderRepositoryImpl) UpdatePaymentAndOrderPaidAmount(ctx context.Context, paymentID uuid.UUID, paymentUpdates map[string]interface{}, orderID uuid.UUID, newPaidAmount int64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).Where("id = ?", paymentID).Updates(paymentUpdates).Error; err != nil {
			return err
		}
		if err := tx.Scopes(database.TenantScope(ctx)).Model(&models.Order{}).Where("id = ?", orderID).Update("paid_amount", newPaidAmount).Error; err != nil {
			return err
		}
		return nil
	})
}

// CancelOrderWithPayments updates order status to DIBATALKAN and rejects related payments atomically.
func (r *OrderRepositoryImpl) CancelOrderWithPayments(ctx context.Context, orderID uuid.UUID, orderUpdates map[string]interface{}, paymentUpdates map[string]interface{}) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Scopes(database.TenantScope(ctx)).Model(&models.Order{}).Where("id = ?", orderID).Updates(orderUpdates).Error; err != nil {
			return err
		}
		if err := tx.Scopes(database.TenantScopeByOrder(ctx)).Model(&models.Payment{}).Where("order_id = ?", orderID).Updates(paymentUpdates).Error; err != nil {
			return err
		}
		return nil
	})
}

