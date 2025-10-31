package service

import (
	"context"
	"errors"

	dommedia "simple-go/internal/domain/media"
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/volume"
	"simple-go/internal/repository"
	"simple-go/pkg/epub/transformer"
	"simple-go/pkg/logger"

	"gorm.io/gorm"
)

type NovelService struct {
	uow                repository.UnitOfWork
	novelRepo          repository.NovelRepository
	mediaSrvc          *MediaService
	volumeSrvc         *VolumeService
	epubSrvc           *EpubService
	transformerFactory *transformer.EpubTransformerFactory
}

func NewNovelService(uow repository.UnitOfWork, novelRepo repository.NovelRepository, mediaSrvc *MediaService, volumeSrvc *VolumeService, epubSrvc *EpubService) *NovelService {
	return &NovelService{
		uow:                uow,
		novelRepo:          novelRepo,
		mediaSrvc:          mediaSrvc,
		volumeSrvc:         volumeSrvc,
		epubSrvc:           epubSrvc,
		transformerFactory: transformer.NewEpubTransformerFactory(),
	}
}

func (s *NovelService) CreateNovelWithTranslation(
	ctx context.Context,
	creatorID string,
	dto novel.CreateNovelDTO,
) (*novel.Novel, *novel.NovelTranslation, error) {
	var newNovel *novel.Novel
	var newTranslation *novel.NovelTranslation

	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		novelToCreate := &novel.Novel{
			CreatedBy:        creatorID,
			OriginalLanguage: dto.OriginalLanguage,
			OriginalAuthor:   dto.OriginalAuthor,
			Source:           dto.Source,
			Status:           dto.Status,
		}

		createdNovel, err := provider.Novel().Create(ctx, novelToCreate)
		if err != nil {
			logger.Error(err, "failed to create novel in database")
			return errors.New("unable to create novel")
		}
		newNovel = createdNovel

		translationToCreate := &novel.NovelTranslation{
			NovelID:     newNovel.ID,
			Lang:        dto.OriginalLanguage,
			Title:       dto.Title,
			Description: dto.Description,
		}

		createdTranslation, err := provider.Novel().CreateTranslation(ctx, translationToCreate)
		if err != nil {
			logger.Error(err, "failed to create novel translation")
			return errors.New("unable to create novel translation")
		}
		newTranslation = createdTranslation

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return newNovel, newTranslation, nil
}

func (s *NovelService) GetByID(ctx context.Context, id, lang string) (*novel.NovelResponseDTO, error) {
	n, err := s.novelRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		logger.Error(err, "failed to get novel by ID")
		return nil, errors.New("unable to retrieve novel")
	}

	res := novel.MapNovelToDTO(*n, lang)
	return &res, nil
}

func (s *NovelService) GetAll(ctx context.Context, limit, offset int, title, lang string) ([]novel.NovelResponseDTO, int64, error) {
	var novels []novel.Novel
	var err error

	if lang != "" {
		novels, err = s.novelRepo.GetAllByLang(ctx, lang, limit, offset)
	} else {
		novels, err = s.novelRepo.GetAll(ctx, limit, offset, title, lang)
	}

	if err != nil {
		logger.Error(err, "failed to get all novels")
		return nil, 0, errors.New("unable to retrieve novels")
	}

	count, err := s.novelRepo.Count(ctx)

	if err != nil {
		logger.Error(err, "failed to count novels")
		return nil, 0, errors.New("unable to retrieve novels")
	}

	response := make([]novel.NovelResponseDTO, len(novels))

	for i, n := range novels {
		response[i] = novel.MapNovelToDTO(n, lang)
	}
	return response, count, nil
}

func (s *NovelService) Delete(ctx context.Context, id string) error {
	if affected, err := s.novelRepo.Delete(ctx, id); err != nil {
		logger.Error(err, "failed to delete novel")
		return errors.New("unable to delete novel")
	} else if affected == 0 {
		return errors.New("novel not found")
	}
	return nil
}

func (s *NovelService) UpdateCoverMedia(
	ctx context.Context,
	id string,
	dto novel.UpdateCoverMediaDTO,
) error {
	uploadParams := dommedia.UploadAndSaveDTO{
		Name:       dto.FileName,
		FileBytes:  dto.FileBytes,
		UploaderID: dto.UploaderID,
	}

	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		savedMedia, _, err := s.mediaSrvc.UploadAndSaveWithRepo(ctx, provider.Media(), uploadParams)
		if err != nil {
			logger.Error(err, "failed to upload and save media for novel cover")
			return errors.New("unable to upload cover media")
		}

		_, err = provider.Novel().UpdateCoverMedia(ctx, id, savedMedia.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("novel not found")
			}
			logger.Error(err, "failed to update novel cover media")
			return errors.New("unable to update novel cover")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *NovelService) CreateTranslation(
	ctx context.Context,
	dto novel.CreateNovelTranslationDTO,
) (*novel.NovelTranslation, error) {
	_, err := s.novelRepo.GetByID(ctx, dto.NovelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		logger.Error(err, "failed to verify novel exists")
		return nil, errors.New("unable to create translation")
	}

	existing, err := s.novelRepo.GetTranslation(ctx, dto.NovelID, dto.Lang)
	if err == nil && existing != nil {
		return nil, errors.New("translation for this language already exists")
	}

	nt := &novel.NovelTranslation{
		NovelID:     dto.NovelID,
		Lang:        dto.Lang,
		Title:       dto.Title,
		Description: dto.Description,
	}

	createdTranslation, err := s.novelRepo.CreateTranslation(ctx, nt)
	if err != nil {
		logger.Error(err, "failed to create novel translation")
		return nil, errors.New("unable to create translation")
	}

	// 5. Return the created entity
	return createdTranslation, nil
}

func (s *NovelService) DeleteTranslation(ctx context.Context, id string) error {
	affected, err := s.novelRepo.DeleteTranslation(ctx, id)
	if err != nil {
		logger.Error(err, "failed to delete novel translation")
		return errors.New("unable to delete translation")
	}

	if affected == 0 {
		return errors.New("novel translation not found")
	}

	return nil
}

func (s *NovelService) GetNovelVolumes(ctx context.Context, novelID, lang string) ([]volume.VolumeResponseDTO, error) {
	return s.volumeSrvc.GetNovelVolumes(ctx, novelID, lang)
}

func (s *NovelService) ProcessEpubUpload(ctx context.Context, fileBytes []byte) (*transformer.EpubProcessResult, error) {
	rawEpub, err := s.epubSrvc.UploadAndExtractRawEpub(ctx, fileBytes)
	if err != nil {
		logger.Error(err, "Failed to extract raw EPUB")
		return nil, err
	}

	// Step 2: Auto-detect source and get appropriate transformer
	tr, err := s.transformerFactory.DetectAndGetTransformer(rawEpub)
	if err != nil {
		logger.Error(err, "Failed to detect EPUB source type")
		return nil, err
	}

	// Step 3: Transform EPUB content to database-ready format
	novelData, err := tr.TransformToNovelData(ctx, rawEpub)
	if err != nil {
		logger.Error(err, "Failed to transform novel data")
		return nil, err
	}

	volumes, err := tr.TransformToVolumes(ctx, rawEpub)
	if err != nil {
		logger.Error(err, "Failed to transform volumes")
		return nil, err
	}

	chapters, err := tr.TransformToChapters(ctx, rawEpub)
	if err != nil {
		logger.Error(err, "Failed to transform chapters")
		return nil, err
	}

	result := &transformer.EpubProcessResult{
		RawContent:    rawEpub,
		NovelData:     novelData,
		Volumes:       volumes,
		Chapters:      chapters,
		SourceType:    tr.GetSourceType(),
		TotalChapters: len(chapters),
		TotalVolumes:  len(volumes),
	}

	return result, nil
}

func (s *NovelService) ProcessAndSaveEpubUpload(ctx context.Context, fileBytes []byte, creatorID string) (*transformer.EpubProcessResult, error) {
	result, err := s.ProcessEpubUpload(ctx, fileBytes)
	if err != nil {
		return nil, err
	}

	err = s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		persist := &epubPersistence{
			ctx:       ctx,
			creatorID: creatorID,
			result:    result,
			provider:  provider,
			mediaSrvc: s.mediaSrvc,
		}
		return persist.run()
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
