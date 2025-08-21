package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"log"
	"novel-server/tools"
	"novel-server/web/models"
)

type NovelController struct {
	Ctx iris.Context
}

// GET /novels
func (c *NovelController) Get() ([]models.Novel, error) {
	novels, ok := tools.MySQLGetAllNovels()
	if !ok || len(novels) == 0 {
		log.Printf("未获取到小说信息")
		return []models.Novel{}, nil
	}
	for _, novel := range novels {
		log.Printf("小说: %v", novel)
		break
	}
	return novels, nil
}

// GET /novels/{id}
func (c *NovelController) GetBy(novel_id string) (models.Novel, error) {
	novel, ok := tools.MySQLGetNovelByID(novel_id)
	if !ok {
		log.Printf("获取小说 ID=%v 失败", novel_id)
		return models.Novel{}, fmt.Errorf("获取小说失败")
	}
	return novel, nil
}

// GET /novels/{id}
func (c *NovelController) GetBy_novelid(novel_id string) (models.Novel, error) {
	novel, ok := tools.MySQLGetNovelByID(novel_id)
	if !ok {
		log.Printf("获取小说 ID=%v 失败", novel_id)
		return models.Novel{}, fmt.Errorf("获取小说失败")
	}
	return novel, nil
}

// POST /novels
func (c *NovelController) Post(newNovel models.Novel) (models.Novel, error) {
	id, err := tools.MySQLCreateNovel(newNovel)
	if err != nil {
		log.Printf("创建小说失败: %v", err)
		return newNovel, err
	}
	newNovel.ID = id
	return newNovel, nil
}

// PUT /novels/{id}
func (c *NovelController) PutBy(novel_id string, updated models.Novel) (models.Novel, error) {
	updated.Novel_Id = novel_id
	ok, err := tools.MySQLUpdateNovel(updated)
	if err != nil || !ok {
		log.Printf("更新小说 ID=%v 失败: %v", novel_id, err)
		return updated, fmt.Errorf("更新小说失败")
	}
	return updated, nil
}

// DELETE /novels/{id}
func (c *NovelController) DeleteBy(novel_id string) tools.Response {
	ok, err := tools.MySQLDeleteNovel(novel_id)
	if err != nil || !ok {
		log.Printf("删除小说 ID=%v 失败: %v", novel_id, err)
		return tools.Fail(tools.ErrorCode.CodeDeleteNovelFailed)
	}
	return tools.Success(map[string]any{
		"novel_id": novel_id,
	})
}

// GET /novels/{novel_id}/chapters
func (c *NovelController) GetByNovelId(novel_id string) ([]models.Chapter, error) {
	log.Printf("获取小说  ")
	chapters, ok := tools.MySQLGetChaptersByNovelID(novel_id)
	if !ok {
		log.Printf("获取小说 ID=%V 的章节失败", novel_id)
		return nil, fmt.Errorf("获取章节失败")
	}
	return chapters, nil
}
