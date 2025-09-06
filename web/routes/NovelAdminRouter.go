package routes

import (
	"novel-server/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterNovelAdminRouter(app *iris.Application) {
	// 小说相关路由
	mvc.New(app.Party("/admin/novels")).Handle(new(controllers.NovelAdminController))
	mvc.New(app.Party("/admin/novels/{novelid:string}/audit")).Handle(new(controllers.NovelAdminController))
	mvc.New(app.Party("/admin/novels/{novelid:string}/chapters")).Handle(new(controllers.Novel_Chaptes_Admin_Controller))
	mvc.New(app.Party("/admin/chapters/{novelid:string}/{chapterIndex:int}")).Handle(new(controllers.ChapterController))
	mvc.New(app.Party("/admin/chapters/{novelid:string}")).Handle(new(controllers.ChapterController))
}
