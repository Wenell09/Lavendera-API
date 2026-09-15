package contextutil

import "context"

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	TenantIDKey contextKey = "tenant_id"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(UserIDKey).(string)
	return value, ok
}

func TenantIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(TenantIDKey).(string)
	return value, ok
}
