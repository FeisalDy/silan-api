package novelgenre

import (
	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/novel"
)

type NovelGenre struct {
	NovelID string      `gorm:"type:uuid;primaryKey"`
	GenreID string      `gorm:"type:uuid;primaryKey"`
	Novel   novel.Novel `gorm:"foreignKey:NovelID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Genre   genre.Genre `gorm:"foreignKey:GenreID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (NovelGenre) TableName() string {
	return "novel_genres"
}
