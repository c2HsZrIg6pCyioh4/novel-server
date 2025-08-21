package routes

import (
	"novel-server/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterNovelRouter(app *iris.Application) {
	// 小说相关路由
	mvc.New(app.Party("/novels")).Handle(new(controllers.NovelController))
	mvc.New(app.Party("/novels/{novelid:string}/chapters")).Handle(new(controllers.Novel_Chaptes_Controller))
	mvc.New(app.Party("/chapters/{novelid:string}/{chapterIndex:int}")).Handle(new(controllers.ChapterController))
	mvc.New(app.Party("/chapters/{novelid:string}")).Handle(new(controllers.ChapterController))
}
