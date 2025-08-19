package routes

import (
	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app *iris.Application) {
	// Novel routes
	RegisterNovelRouter(app)
}
