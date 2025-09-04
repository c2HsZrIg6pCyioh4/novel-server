package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"novel-server/web/controllers"
)

func RegisterImageRouter(app *iris.Application) {
	// 健康检查接口
	mvc.New(app.Party("/images/uploads")).Handle(new(controllers.UploadController))
	mvc.New(app.Party("/images/{format:string}/{year:string}/{month:string}/{day:string}/{filename:string}")).Handle(new(controllers.ImageController))

}
