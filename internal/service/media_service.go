package service

import (
	"context"
	"errors"
	"strings"

	dommedia "simple-go/internal/domain/media"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"
	mediapkg "simple-go/pkg/media"
)

type MediaService struct {
	mediaRepo     repository.MediaRepository
	uploadService *UploadService
}

func NewMediaService(mediaRepo repository.MediaRepository, uploadService *UploadService) *MediaService {
	return &MediaService{mediaRepo: mediaRepo, uploadService: uploadService}
}

func (s *MediaService) UploadAndSaveMedia(ctx context.Context, p dommedia.UploadAndSaveDTO) (*dommedia.Media, *mediapkg.Response, error) {
	sources := 0
	if len(p.FileBytes) > 0 {
		sources++
	}
	if p.Base64 != "" {
		sources++
	}
	if p.URL != "" {
		sources++
	}
	if sources == 0 {
		return nil, nil, errors.New("no media source provided: file bytes, base64, or URL required")
	}
	if sources > 1 {
		return nil, nil, errors.New("multiple media sources provided: provide exactly one of file bytes, base64, or URL")
	}

	name := strings.TrimSpace(p.Name)
	if name == "" {
		name = "upload"
	}

	var (
		uploadResp mediapkg.Response
		err        error
	)
	switch {
	case len(p.FileBytes) > 0:
		uploadResp, err = s.uploadService.UploadFile(ctx, name, p.FileBytes, p.TTL)
	case p.Base64 != "":
		uploadResp, err = s.uploadService.UploadBase64(ctx, name, p.Base64, p.TTL)
	case p.URL != "":
		uploadResp, err = s.uploadService.UploadURL(ctx, name, p.URL, p.TTL)
	}

	if err != nil {
		logger.Error(err, "Upload Failed")
		return nil, nil, errors.New("failed to upload the media")
	}

	url := uploadResp.Data.URL
	var mime *string
	if uploadResp.Data.Image.Mime != "" {
		m := uploadResp.Data.Image.Mime
		mime = &m
	} else if p.MimeType != nil {
		mime = p.MimeType
	}

	var sizePtr *int64
	if uploadResp.Data.Size > 0 {
		sz := int64(uploadResp.Data.Size)
		sizePtr = &sz
	} else if p.FileSize != nil {
		sizePtr = p.FileSize
	}

	var typ *string
	if mime != nil {
		parts := strings.SplitN(*mime, "/", 2)
		if len(parts) > 0 && parts[0] != "" {
			t := parts[0]
			typ = &t
		}
	}

	newMedia := &dommedia.Media{
		URL:         &url,
		Type:        typ,
		Description: p.Description,
		FileSize:    sizePtr,
		MimeType:    mime,
	}

	if p.UploaderID != "" {
		uid := p.UploaderID
		newMedia.UploadedBy = &uid
	}

	savedMedia, err := s.mediaRepo.Create(ctx, newMedia)
	if err != nil {
		logger.Error(err, "failed to save media metadata")
		return nil, nil, errors.New("failed to save media metadata")
	}

	return savedMedia, &uploadResp, nil
}

// UploadAndSaveWithRepo performs the upload and saves metadata using the provided repository.
// This is useful to participate in a UnitOfWork transaction where the repo comes from the provider.
func (s *MediaService) UploadAndSaveWithRepo(ctx context.Context, repo repository.MediaRepository, p dommedia.UploadAndSaveDTO) (*dommedia.Media, *mediapkg.Response, error) {
	// Reuse the same logic but don't use s.mediaRepo; use the provided repo instead.
	sources := 0
	if len(p.FileBytes) > 0 {
		sources++
	}
	if p.Base64 != "" {
		sources++
	}
	if p.URL != "" {
		sources++
	}
	if sources == 0 {
		return nil, nil, errors.New("no media source provided: file bytes, base64, or URL required")
	}
	if sources > 1 {
		return nil, nil, errors.New("multiple media sources provided: provide exactly one of file bytes, base64, or URL")
	}

	name := strings.TrimSpace(p.Name)
	if name == "" {
		name = "upload"
	}

	var (
		uploadResp mediapkg.Response
		err        error
	)
	switch {
	case len(p.FileBytes) > 0:
		uploadResp, err = s.uploadService.UploadFile(ctx, name, p.FileBytes, p.TTL)
	case p.Base64 != "":
		uploadResp, err = s.uploadService.UploadBase64(ctx, name, p.Base64, p.TTL)
	case p.URL != "":
		uploadResp, err = s.uploadService.UploadURL(ctx, name, p.URL, p.TTL)
	}
	if err != nil {
		logger.Error(err, "Upload Failed")
		return nil, nil, errors.New("failed to upload the media")
	}

	url := uploadResp.Data.URL
	var mime *string
	if uploadResp.Data.Image.Mime != "" {
		m := uploadResp.Data.Image.Mime
		mime = &m
	} else if p.MimeType != nil {
		mime = p.MimeType
	}

	var sizePtr *int64
	if uploadResp.Data.Size > 0 {
		sz := int64(uploadResp.Data.Size)
		sizePtr = &sz
	} else if p.FileSize != nil {
		sizePtr = p.FileSize
	}

	var typ *string
	if mime != nil {
		parts := strings.SplitN(*mime, "/", 2)
		if len(parts) > 0 && parts[0] != "" {
			t := parts[0]
			typ = &t
		}
	}

	newMedia := &dommedia.Media{
		URL:         &url,
		Type:        typ,
		Description: p.Description,
		FileSize:    sizePtr,
		MimeType:    mime,
	}
	if p.UploaderID != "" {
		uid := p.UploaderID
		newMedia.UploadedBy = &uid
	}

	savedMedia, err := repo.Create(ctx, newMedia)
	if err != nil {
		logger.Error(err, "failed to save media metadata")
		return nil, nil, errors.New("failed to save media metadata")
	}

	return savedMedia, &uploadResp, nil
}
