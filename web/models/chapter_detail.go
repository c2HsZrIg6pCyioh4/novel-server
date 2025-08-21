package models

// ChapterDetail 表示小说的单章节详情
type ChapterDetail struct {
	Title        string `json:"title" gorm:"size:255;not null"` // 章节标题
	ChapterIndex int    `json:"chapter_index"`                  // 章节顺序索引
}
