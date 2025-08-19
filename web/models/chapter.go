package models

import "time"

type Chapter struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	NovelID      int64     `json:"novel_id" gorm:"index;not null"`
	Title        string    `json:"title" gorm:"size:255;not null"`
	Content      string    `json:"content" gorm:"type:longtext"`
	WordCount    int       `json:"word_count"`
	ChapterIndex int       `json:"chapter_index"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
