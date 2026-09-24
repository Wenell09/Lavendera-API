package appcontext

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	TenantIDKey contextKey = "tenant_id"
	roleKey     contextKey = "role"
)

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return value, ok
}

func TenantIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(TenantIDKey).(uuid.UUID)
	return value, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
