package storage

import (
	"context"
	"io"

	"github.com/Wenell09/lavendera-api/internal/shared/config"
	"github.com/google/uuid"
)

type StorageRepository interface {
	GetPublicURL(objectPath string) string
	Upload(ctx context.Context, tenantID uuid.UUID, file io.Reader, contentType string, extension string) (string, error)
	Delete(ctx context.Context, objectPath string) error
}

func ProvideStorageRepository() StorageRepository {
	return NewStorageRepository(
		config.ENV.SupabaseURL,
		config.ENV.SupabaseKey,
		config.ENV.SupabaseBucket,
	)
}
