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

// GET /chapters/{novel_id}/{chapter_index}
// 根据小说ID + 章节索引获取章节
func (c *ChapterController) Get(novel_id string, chapterIndex int) (models.Chapter, error) {
	chapter, ok := tools.MySQLGetChapterByNovelIDAndIndex(novel_id, chapterIndex)
	if !ok {
		log.Printf("获取小说 ID=%v 第 %d 章失败", novel_id, chapterIndex)
		return models.Chapter{}, fmt.Errorf("获取章节失败")
	}
	return chapter, nil
}

// POST /chapters/{novel_id}
// 新增章节
func (c *ChapterController) Post(novel_id string, newChapter models.Chapter) (models.Chapter, error) {
	newChapter.NovelID = novel_id
	// 查询该小说最新的章节序号
	latestIndex, err := tools.MySQLGetLatestChapterIndex(novel_id)
	if err != nil {
		log.Printf("获取小说 ID=%V 的最新章节序号失败: %v", novel_id, err)
		return newChapter, err
	}

	// 新章节序号 = 最新序号 + 1
	newChapter.ChapterIndex = latestIndex + 1
	id, err := tools.MySQLCreateChapter(newChapter)
	if err != nil {
		log.Printf("创建小说 ID=%v 的章节失败: %v", novel_id, err)
		return newChapter, err
	}
	newChapter.ID = id

	// 2️⃣ 同步插入章节目录表
	chapterDetail := models.ChapterDetail{
		Title:        newChapter.Title,
		ChapterIndex: newChapter.ChapterIndex,
	}
	err = tools.MySQLCreateChapterDetail(newChapter.NovelID, chapterDetail)
	if err != nil {
		log.Printf("同步更新小说 ID=%v 章节目录失败: %v", novel_id, err)
		// 这里可以选择忽略错误或者返回错误，看业务需求
	}
	return newChapter, nil
}

// PUT /chapters/{novel_id}/{chapter_index}
// 更新章节
func (c *ChapterController) Put(novel_id string, chapterIndex int, updated models.Chapter) (models.Chapter, error) {
	ok, err := tools.MySQLUpdateChapterByNovelIDAndIndex(novel_id, chapterIndex, updated)
	if err != nil || !ok {
		log.Printf("更新小说 ID=%v 第 %d 章失败: %v", novel_id, chapterIndex, err)
		return updated, fmt.Errorf("更新章节失败")
	}
	return updated, nil
}

// DELETE /chapters/{novel_id}/{chapter_index}
// 删除章节（只允许删除最新章节）
func (c *ChapterController) Delete(novel_id string, chapterIndex int) tools.Response {
	// 1. 查询该小说最新章节的序号
	latestIndex, err := tools.MySQLGetLatestChapterIndex(novel_id)
	if err != nil {
		log.Printf("获取小说 ID=%V 的最新章节序号失败: %v", novel_id, err)
		return tools.Fail(tools.ErrorCode.CodeGetLatestIndexFail)
	}

	// 2. 判断是否是最新章节
	if chapterIndex != latestIndex {
		log.Printf("小说 ID=%v 尝试删除非最新章节 %d，最新章节是 %d", novel_id, chapterIndex, latestIndex)
		return tools.Fail(tools.ErrorCode.CodeNotLatestChapter)
	}

	// 3. 删除最新章节
	ok, err := tools.MySQLDeleteChapterByNovelIDAndIndex(novel_id, chapterIndex)
	if err != nil || !ok {
		log.Printf("删除小说 ID=%v 第 %d 章失败: %v", novel_id, chapterIndex, err)
		return tools.Fail(tools.ErrorCode.CodeDeleteChapterFail)
	}

	return tools.Success(map[string]any{
		"novel_id":      novel_id,
		"chapter_index": chapterIndex,
		"deleted":       true,
	})
}
