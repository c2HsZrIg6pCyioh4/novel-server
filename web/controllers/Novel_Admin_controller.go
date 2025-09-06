package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"log"
	"novel-server/tools"
	"novel-server/web/models"
	"time"
)

type NovelAdminController struct {
	Ctx iris.Context
}

// GET /novels
func (c *NovelAdminController) Get() ([]models.Novel, error) {
	novels, ok := tools.MySQLGetAllNovelsForAudit()
	log.Printf("1111")
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
func (c *NovelAdminController) GetBy(novel_id string) (models.Novel, error) {
	log.Printf("22222")
	novel, ok := tools.MySQLGetNovelByIDForAudit(novel_id)
	if !ok {
		log.Printf("获取小说 ID=%v 失败", novel_id)
		return models.Novel{}, fmt.Errorf("获取小说失败")
	}
	return novel, nil
}

// GET /novels/{id}
func (c *NovelAdminController) GetBy_novelid(novel_id string) (models.Novel, error) {
	log.Printf("33333")
	novel, ok := tools.MySQLGetNovelByIDForAudit(novel_id)
	if !ok {
		log.Printf("获取小说 ID=%v 失败", novel_id)
		return models.Novel{}, fmt.Errorf("获取小说失败")
	}
	return novel, nil
}

// POST /admin/novels/{novel_id}/audit
func (c *NovelAdminController) PostByAudit(novel_id string) tools.Response {
	log.Printf("执行了这个post修改小说审核状态")
	var body struct {
		Status int `json:"status"`
	}

	if err := c.Ctx.ReadJSON(&body); err != nil {
		log.Printf("解析请求体失败: %v", err)
		return tools.Fail(tools.ErrorCode.CodeDeleteChapterFail)
	}

	// 获取当前登录用户（假设你有认证中间件）
	// 如果没有用户系统，可以直接写死一个管理员名
	auditBy := "admin" // 示例，应该从 session 或 token 中取
	now := time.Now()
	ok, err := tools.MySQLUpdateNovelAudit(novel_id, body.Status, auditBy, now)
	if err != nil || !ok {
		log.Printf("审核更新失败 novel_id=%v err=%v", novel_id, err)
		return tools.Fail(tools.ErrorCode.CodeNotLatestChapter)
	}

	return tools.Success(map[string]any{
		"novel_id":     novel_id,
		"audit_status": body.Status,
		"audit_by":     auditBy,
		"audit_at":     now,
	})
}

// POST /novels
func (c *NovelAdminController) Post(newNovel models.Novel) (models.Novel, error) {
	log.Printf("执行了这个post创建小说")
	id, err := tools.MySQLCreateNovel(newNovel)
	if err != nil {
		log.Printf("创建小说失败: %v", err)
		return newNovel, err
	}
	newNovel.ID = id
	return newNovel, nil
}

// PUT /novels/{id}
func (c *NovelAdminController) PutBy(novel_id string, updated models.Novel) (models.Novel, error) {
	updated.Novel_Id = novel_id
	ok, err := tools.MySQLUpdateNovel(updated)
	if err != nil || !ok {
		log.Printf("更新小说 ID=%v 失败: %v", novel_id, err)
		return updated, fmt.Errorf("更新小说失败")
	}
	return updated, nil
}

// DELETE /novels/{id}
func (c *NovelAdminController) DeleteBy(novel_id string) tools.Response {
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
func (c *NovelAdminController) GetByNovelId(novel_id string) ([]models.Chapter, error) {
	log.Printf("获取小说  ")
	chapters, ok := tools.MySQLGetChaptersByNovelID(novel_id)
	if !ok {
		log.Printf("获取小说 ID=%V 的章节失败", novel_id)
		return nil, fmt.Errorf("获取章节失败")
	}
	return chapters, nil
}
