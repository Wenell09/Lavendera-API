package utils

import (
	"context"

	"github.com/Wenell09/lavendera-api/internal/shared/appcontext"
	"github.com/sirupsen/logrus"
)

// WithContext membuat logrus Entry baru yang sudah otomatis menyisipkan tenant_id jika ada di context
func LogWithContext(baseLogger *logrus.Logger, ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(baseLogger)
	if tenantID, ok := appcontext.TenantIDFromContext(ctx); ok {
		entry = entry.WithField("tenant_id", tenantID)
	}
	return entry
}
