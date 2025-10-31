package novelgenre

import (
	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/novel"
)

type NovelGenre struct {
	NovelID string      `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	GenreID string      `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Novel   novel.Novel `gorm:"foreignKey:NovelID;references:ID"`
	Genre   genre.Genre `gorm:"foreignKey:GenreID;references:ID"`
}

func (NovelGenre) TableName() string {
	return "novel_genres"
}
