package controllers

import (
	"novel-server/tools"

	"github.com/kataras/iris/v12"
)

// HealthController 健康检查
type HealthController struct {
	Ctx iris.Context
}

// GET /health
func (c *HealthController) Get() tools.Response {
	return tools.Success(map[string]string{
		"services": "ok",
	})
}
