package appcontext

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	TenantIDKey contextKey = "tenant_id"
)

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return value, ok
}

func TenantIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(TenantIDKey).(uuid.UUID)
	return value, ok
}
