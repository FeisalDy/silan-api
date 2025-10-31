package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domchapter "simple-go/internal/domain/chapter"
	dommedia "simple-go/internal/domain/media"
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/volume"
	"simple-go/internal/repository"
	"simple-go/pkg/epub/transformer"
	"simple-go/pkg/logger"
)

type epubPersistence struct {
	ctx       context.Context
	creatorID string
	result    *transformer.EpubProcessResult
	provider  repository.RepositoryProvider
	mediaSrvc *MediaService

	novel   *novel.Novel
	volumes []*volume.Volume
}

func (p *epubPersistence) run() error {
	coverMediaID, err := p.uploadCoverImage()
	if err != nil {
		return err
	}

	if err := p.createNovel(coverMediaID); err != nil {
		return err
	}

	if err := p.createNovelTranslation(); err != nil {
		return err
	}

	if err := p.syncTags(); err != nil {
		return err
	}

	if err := p.createVolumes(); err != nil {
		return err
	}

	return p.createChapters()
}

func (p *epubPersistence) uploadCoverImage() (*string, error) {
	if p.mediaSrvc == nil || len(p.result.NovelData.CoverImage) == 0 {
		return nil, nil
	}

	uploadParams := dommedia.UploadAndSaveDTO{
		Name:       fmt.Sprintf("%s-cover", p.result.NovelData.Title),
		FileBytes:  p.result.NovelData.CoverImage,
		UploaderID: p.creatorID,
	}

	savedMedia, _, err := p.mediaSrvc.UploadAndSaveWithRepo(p.ctx, p.provider.Media(), uploadParams)
	if err != nil {
		logger.Error(err, "Failed to upload cover image, continuing without cover")
		return nil, nil
	}

	return &savedMedia.ID, nil
}

func (p *epubPersistence) createNovel(coverMediaID *string) error {
	newNovel := &novel.Novel{
		CreatedBy:        p.creatorID,
		OriginalLanguage: p.result.NovelData.OriginalLanguage,
		OriginalAuthor:   optionalStringPtr(p.result.NovelData.OriginalAuthor),
		CoverMediaID:     coverMediaID,
	}

	createdNovel, err := p.provider.Novel().Create(p.ctx, newNovel)
	if err != nil {
		logger.Error(err, "Failed to create novel")
		return errors.New("failed to create novel")
	}

	p.novel = createdNovel
	return nil
}

func (p *epubPersistence) createNovelTranslation() error {
	if p.novel == nil {
		return errors.New("novel must be created before creating translations")
	}

	description := p.result.NovelData.Description

	novelTranslation := &novel.NovelTranslation{
		NovelID:     p.novel.ID,
		Lang:        p.result.NovelData.OriginalLanguage,
		Title:       p.result.NovelData.Title,
		Description: optionalStringPtr(description),
	}

	if _, err := p.provider.Novel().CreateTranslation(p.ctx, novelTranslation); err != nil {
		logger.Error(err, "Failed to create novel translation")
		return errors.New("failed to create novel translation")
	}

	return nil
}

func (p *epubPersistence) syncTags() error {
	if len(p.result.NovelData.Tags) == 0 {
		return nil
	}

	tags, err := p.provider.Tag().FindOrCreateByNames(p.ctx, p.result.NovelData.Tags)
	if err != nil {
		logger.Error(err, "Failed to find/create tags")
		return errors.New("failed to process tags")
	}

	var tagIDs []string
	for _, t := range tags {
		tagIDs = append(tagIDs, t.ID)
	}

	if err := p.provider.NovelTag().LinkTagsToNovel(p.ctx, p.novel.ID, tagIDs); err != nil {
		logger.Error(err, "Failed to link tags to novel")
		return errors.New("failed to link tags")
	}

	return nil
}

func (p *epubPersistence) createVolumes() error {
	volumesData := p.result.Volumes
	if len(volumesData) == 0 {
		volumesData = []transformer.VolumeData{{
			Number:    1,
			Title:     "Volume 1",
			IsVirtual: true,
		}}
	}

	p.volumes = make([]*volume.Volume, len(volumesData))
	for i, volData := range volumesData {
		newVolume := &volume.Volume{
			Number:           defaultVolumeNumber(volData.Number, i+1),
			OriginalLanguage: p.result.NovelData.OriginalLanguage,
			NovelID:          p.novel.ID,
			IsVirtual:        volData.IsVirtual,
		}

		createdVolume, err := p.provider.Volume().Create(p.ctx, newVolume)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Failed to create volume %d", newVolume.Number))
			return fmt.Errorf("failed to create volume %d", newVolume.Number)
		}

		p.volumes[i] = createdVolume
	}

	return nil
}

func (p *epubPersistence) createChapters() error {
	if len(p.volumes) == 0 {
		return errors.New("volumes must be created before chapters")
	}

	createdCount := 0
	for _, chapterData := range p.result.Chapters {
		volume := p.resolveVolumeForChapter(chapterData.VolumeIndex)
		if volume == nil {
			logger.Error(nil, fmt.Sprintf("Skipping chapter due to invalid volume index %d", chapterData.VolumeIndex))
			continue
		}

		newChapter := &domchapter.Chapter{
			Number:   chapterData.OrderNum,
			VolumeID: volume.ID,
		}

		createdChapter, err := p.provider.Chapter().Create(p.ctx, newChapter)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Failed to create chapter %d", chapterData.OrderNum))
			return fmt.Errorf("failed to create chapter %d", chapterData.OrderNum)
		}

		chapterTranslation := &domchapter.ChapterTranslation{
			ChapterID: createdChapter.ID,
			Lang:      p.result.NovelData.OriginalLanguage,
			Title:     chapterData.Title,
			Content:   chapterData.Content,
		}

		if _, err := p.provider.Chapter().CreateTranslation(p.ctx, chapterTranslation); err != nil {
			logger.Error(err, fmt.Sprintf("Failed to create chapter %d translation", chapterData.OrderNum))
			return errors.New("failed to create chapter translation")
		}

		createdCount++
	}

	return nil
}

func (p *epubPersistence) resolveVolumeForChapter(volumeIndex int) *volume.Volume {
	if volumeIndex >= 0 && volumeIndex < len(p.volumes) {
		return p.volumes[volumeIndex]
	}
	if len(p.volumes) == 0 {
		return nil
	}
	return p.volumes[0]
}

func optionalStringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	result := trimmed
	return &result
}

func defaultVolumeNumber(number int, fallback int) int {
	if number > 0 {
		return number
	}
	return fallback
}
