package routes

import (
	"novel-server/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterNovelRouter(app *iris.Application) {
	// 小说相关路由
	mvc.New(app.Party("/novels")).Handle(new(controllers.NovelController))
	mvc.New(app.Party("/novels/{id:uint64}/chapters")).Handle(new(controllers.Novel_Chaptes_Controller))
	mvc.New(app.Party("/chapters/{novelID:int64}/{chapterIndex:int}")).Handle(new(controllers.ChapterController))
	mvc.New(app.Party("/chapters/{novelID:int64}")).Handle(new(controllers.ChapterController))
}
