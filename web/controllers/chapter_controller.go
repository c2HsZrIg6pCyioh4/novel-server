package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"log"
	"novel-server/tools"
	"novel-server/web/models"
)

type ChapterController struct {
	Ctx iris.Context
}

// GET /novels/{novel_id}/chapters
func (c *ChapterController) GetByNovelId(novelID int64) ([]models.Chapter, error) {
	chapters, ok := tools.MySQLGetChaptersByNovelID(novelID)
	if !ok {
		log.Printf("获取小说 ID=%d 的章节失败", novelID)
		return nil, fmt.Errorf("获取章节失败")
	}
	return chapters, nil
}

// GET /chapters/{id}
func (c *ChapterController) GetBy(id int64) (models.Chapter, error) {
	println("获取章节")
	chapter, ok := tools.MySQLGetChapterByID(id)
	if !ok {
		log.Printf("获取章节 ID=%d 失败", id)
		return models.Chapter{}, fmt.Errorf("获取章节失败")
	}
	return chapter, nil
}

// POST /novels/{novel_id}/chapters
func (c *ChapterController) PostByNovelId(novelID int64, newChapter models.Chapter) (models.Chapter, error) {
	newChapter.NovelID = novelID
	id, err := tools.MySQLCreateChapter(newChapter)
	if err != nil {
		log.Printf("创建小说 ID=%d 的章节失败: %v", novelID, err)
		return newChapter, err
	}
	newChapter.ID = id
	return newChapter, nil
}

// PUT /chapters/{id}
func (c *ChapterController) PutBy(id int64, updated models.Chapter) (models.Chapter, error) {
	updated.ID = id
	ok, err := tools.MySQLUpdateChapter(updated)
	if err != nil || !ok {
		log.Printf("更新章节 ID=%d 失败: %v", id, err)
		return updated, fmt.Errorf("更新章节失败")
	}
	return updated, nil
}

// DELETE /chapters/{id}
func (c *ChapterController) DeleteBy(id int64) error {
	ok, err := tools.MySQLDeleteChapter(id)
	if err != nil || !ok {
		log.Printf("删除章节 ID=%d 失败: %v", id, err)
		return fmt.Errorf("删除章节失败")
	}
	return nil
}
