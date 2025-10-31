package transformer

import (
	"errors"
	"fmt"
	"simple-go/pkg/epub"
	"simple-go/pkg/logger"
)

type EpubTransformerFactory struct {
	transformers []EpubTransformer
}

func NewEpubTransformerFactory() *EpubTransformerFactory {
	return &EpubTransformerFactory{
		transformers: []EpubTransformer{
			NewSource404NovelDownloaderTransformer(),
			NewSourceDipubdLightnovelCrawlerTransformer(),
		},
	}
}

// DetectAndGetTransformer detects the EPUB source and returns appropriate transformer
func (f *EpubTransformerFactory) DetectAndGetTransformer(content *epub.RawEpub) (EpubTransformer, error) {
	for _, transformer := range f.transformers {
		if transformer.DetectSource(content) {
			logger.Info(fmt.Sprintf("Using transformer: %s", transformer.GetSourceType()))
			return transformer, nil
		}
	}

	return nil, errors.New("unsupported EPUB source format: no compatible transformer found")
}

// GetTransformerByType returns a transformer for a specific source type
func (f *EpubTransformerFactory) GetTransformerByType(sourceType EpubSourceType) (EpubTransformer, error) {
	for _, transformer := range f.transformers {
		if transformer.GetSourceType() == sourceType {
			return transformer, nil
		}
	}

	return nil, fmt.Errorf("transformer not found for source type: %s", sourceType)
}

// RegisterTransformer allows dynamic registration of new transformers
func (f *EpubTransformerFactory) RegisterTransformer(transformer EpubTransformer) {
	f.transformers = append(f.transformers, transformer)
	logger.Info(fmt.Sprintf("Registered new transformer: %s", transformer.GetSourceType()))
}
