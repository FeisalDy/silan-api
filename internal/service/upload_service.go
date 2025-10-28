package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	mediapkg "simple-go/pkg/media"
)

type UploadService struct {
	client     *mediapkg.Client
	defaultTTL uint64
}

func NewUploadService(httpClient *http.Client, apiKey string, defaultTTL uint64) *UploadService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &UploadService{
		client:     mediapkg.NewClient(httpClient, apiKey),
		defaultTTL: defaultTTL,
	}
}

func (s *UploadService) UploadFile(ctx context.Context, name string, file []byte, ttl uint64) (mediapkg.Response, error) {
	if len(file) == 0 {
		return mediapkg.Response{}, errors.New("file is empty")
	}
	if name == "" {
		name = "upload"
	}
	if ttl == 0 {
		ttl = s.defaultTTL
	}
	img, _ := mediapkg.NewImageFromFile(name, ttl, file)
	return s.client.Upload(ctx, img)
}

func (s *UploadService) UploadBase64(ctx context.Context, name string, base64 string, ttl uint64) (mediapkg.Response, error) {
	if base64 == "" {
		return mediapkg.Response{}, errors.New("base64 source is empty")
	}
	if name == "" {
		name = "upload"
	}
	if ttl == 0 {
		ttl = s.defaultTTL
	}
	img, _ := mediapkg.NewImage(name, ttl, base64)
	return s.client.Upload(ctx, img)
}

func (s *UploadService) UploadURL(ctx context.Context, name string, imageURL string, ttl uint64) (mediapkg.Response, error) {
	if imageURL == "" {
		return mediapkg.Response{}, errors.New("image URL is empty")
	}
	if name == "" {
		name = "upload"
	}
	if ttl == 0 {
		ttl = s.defaultTTL
	}
	img, _ := mediapkg.NewImage(name, ttl, imageURL)
	return s.client.Upload(ctx, img)
}
