package noveltag

import (
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/tag"
	"time"
)

type NovelTag struct {
	NovelID   string      `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TagID     string      `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Novel     novel.Novel `gorm:"foreignKey:NovelID;references:ID"`
	Tag       tag.Tag     `gorm:"foreignKey:TagID;references:ID"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
}

func (NovelTag) TableName() string {
	return "novel_tags"
}
