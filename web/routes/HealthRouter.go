package routes

import (
	"novel-server/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterHealthRouter(app *iris.Application) {
	// 健康检查接口
	mvc.New(app.Party("/health")).Handle(new(controllers.HealthController))
}
