package controllers

import (
	"novel-server/tools"

	"github.com/kataras/iris/v12"
)

// HealthController 健康检查
type Auth_Controller struct {
	Ctx iris.Context
}

// Post /validate-token
func (c *Auth_Controller) Post() tools.Response {
	return tools.Success(map[string]string{
		"token": "ok",
	})
}
