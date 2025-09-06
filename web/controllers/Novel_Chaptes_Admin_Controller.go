package controllers

import (
	"github.com/kataras/iris/v12"
	"log"
	"novel-server/tools"
	"novel-server/web/models"
)

type Novel_Chaptes_Admin_Controller struct {
	Ctx iris.Context
}

// GET /admin/novels/{id:uint}/chapters
func (c *Novel_Chaptes_Admin_Controller) Get(novel_id string) ([]models.ChapterDetail, error) {
	chaptersdetail, ok := tools.MySQLGetChaptersDetailByNovelIDForAudit(novel_id)
	if !ok || len(chaptersdetail) == 0 {
		log.Printf("获取小说 ID=%v 的章节为空或失败", novel_id)
		// 返回默认章节
		return []models.ChapterDetail{
			{
				Title:        "暂无章节",
				ChapterIndex: 0,
			},
		}, nil
	}
	return chaptersdetail, nil
}
