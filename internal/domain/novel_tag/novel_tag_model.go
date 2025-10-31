package noveltag

import (
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/tag"
)

type NovelTag struct {
	NovelID string `gorm:"type:uuid;primaryKey"`
	TagID   string `gorm:"type:uuid;primaryKey"`

	Novel novel.Novel `gorm:"foreignKey:NovelID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Tag   tag.Tag     `gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (NovelTag) TableName() string {
	return "novel_tags"
}
