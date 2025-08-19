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
	println("获取小说列表")
	novels, ok := tools.MySQLGetAllNovels()
	if !ok {
		log.Printf("未获取到小说信息")
		return nil, fmt.Errorf("未获取到小说信息")
	}
	for _, novel := range novels {
		log.Printf("小说: %v", novel)
		break
	}
	return novels, nil
}

// GET /novels/{id}
func (c *NovelController) GetBy(id int64) (models.Novel, error) {
	novel, ok := tools.MySQLGetNovelByID(id)
	if !ok {
		log.Printf("获取小说 ID=%d 失败", id)
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
func (c *NovelController) PutBy(id int64, updated models.Novel) (models.Novel, error) {
	updated.ID = id
	ok, err := tools.MySQLUpdateNovel(updated)
	if err != nil || !ok {
		log.Printf("更新小说 ID=%d 失败: %v", id, err)
		return updated, fmt.Errorf("更新小说失败")
	}
	return updated, nil
}

// DELETE /novels/{id}
func (c *NovelController) DeleteBy(id int64) error {
	ok, err := tools.MySQLDeleteNovel(id)
	if err != nil || !ok {
		log.Printf("删除小说 ID=%d 失败: %v", id, err)
		return fmt.Errorf("删除小说失败")
	}
	return nil
}
