package routes

import (
	"novel-server/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterAuthRouter(app *iris.Application) {
	// 健康检查接口
	mvc.New(app.Party("/validate-token")).Handle(new(controllers.Auth_Controller))
	mvc.New(app.Party("/oauth/{provider:string}/config")).Handle(new(controllers.Auth_Config_Controller))
	mvc.New(app.Party("/oauth/{provider:string}/post_callback")).Handle(new(controllers.Auth_Post_Callback_Controller))
	mvc.New(app.Party("/oauth/{provider:string}/callback")).Handle(new(controllers.Auth_Callback_Controller))
}
