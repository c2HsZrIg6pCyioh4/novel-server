package models

import "time"

type Novel struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Novel_Id    string    `json:"novel_id" gorm:"size:255;not null"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Author      string    `json:"author" gorm:"size:100;not null"`
	Category    string    `json:"category" gorm:"size:50"`
	Status      int       `json:"status" gorm:"default:0"` // 0=连载中 1=已完结
	Description string    `json:"description" gorm:"type:text"`
	CoverURL    string    `json:"cover_url" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
