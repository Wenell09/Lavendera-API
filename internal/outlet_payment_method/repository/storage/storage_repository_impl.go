package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type StorageRepositoryImpl struct {
	baseURL    string
	apiKey     string
	bucket     string
	httpClient *http.Client
}

func NewStorageRepository(
	baseURL string,
	apiKey string,
	bucket string,
) StorageRepository {
	return &StorageRepositoryImpl{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		bucket:  bucket,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *StorageRepositoryImpl) GetPublicURL(objectPath string) string {
	return fmt.Sprintf(
		"%s/storage/v1/object/public/%s/%s",
		strings.TrimRight(s.baseURL, "/"),
		url.PathEscape(s.bucket),
		encodeObjectPath(objectPath),
	)
}
func (s *StorageRepositoryImpl) Upload(ctx context.Context, tenantID uuid.UUID, file io.Reader, contentType string, extension string,
) (string, error) {
	filename := fmt.Sprintf(
		"%s/%d%s",
		tenantID.String(),
		time.Now().UnixNano(),
		extension,
	)

	endpoint := fmt.Sprintf(
		"%s/storage/v1/object/%s/%s",
		strings.TrimRight(s.baseURL, "/"),
		url.PathEscape(s.bucket),
		encodeObjectPath(filename),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		file,
	)
	if err != nil {
		return "", fmt.Errorf("create upload request: %w", err)
	}

	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

		return "", fmt.Errorf(
			"supabase upload failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	return filename, nil
}

func (s *StorageRepositoryImpl) Delete(ctx context.Context, objectPath string,
) error {
	endpoint := fmt.Sprintf(
		"%s/storage/v1/object/%s",
		strings.TrimRight(s.baseURL, "/"),
		url.PathEscape(s.bucket),
	)

	payload, err := json.Marshal(map[string][]string{
		"prefixes": {objectPath},
	})
	if err != nil {
		return fmt.Errorf("marshal delete payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		endpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("create delete request: %w", err)
	}

	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

		return fmt.Errorf(
			"supabase delete failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	return nil
}
func encodeObjectPath(objectPath string) string {
	parts := strings.Split(objectPath, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}
