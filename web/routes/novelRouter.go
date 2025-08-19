package routes

import (
	"novel-server/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterNovelRouter(app *iris.Application) {
	// 小说相关路由
	mvc.New(app.Party("/novels")).Handle(new(controllers.NovelController))

	// 章节相关路由
	mvc.New(app.Party("/chapters")).Handle(new(controllers.ChapterController))
}
