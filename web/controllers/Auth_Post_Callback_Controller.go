package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"novel-server/tools"
	"novel-server/web/models"
	"strings"
)

type Auth_Post_Callback_Controller struct {
	Ctx iris.Context
}

// POST /oauth/{provider:string}/callback
func (c *Auth_Post_Callback_Controller) Post(provider string) models.OAuthToken {
	config, err := tools.GetAppConfig("config.yaml")
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString("读取配置失败")
		return models.OAuthToken{}
	}

	switch provider {
	case "apple":
		// 读取回调参数
		code := c.Ctx.FormValue("code")
		state := c.Ctx.FormValue("state")

		if code == "" || state == "" {
			c.Ctx.StatusCode(iris.StatusBadRequest)
			c.Ctx.WriteString("缺少 code 或 state 参数")
			return models.OAuthToken{}
		}

		// 获取 Apple 配置
		appleConfig, ok := config.OAuth["apple"]
		if !ok || appleConfig.RedirectURI == "" {
			c.Ctx.StatusCode(iris.StatusInternalServerError)
			c.Ctx.WriteString("Apple OAuth 配置错误")
			return models.OAuthToken{}
		}

		// 将 post_callback 替换为 callback
		redirectURI := strings.Replace(appleConfig.RedirectURI, "post_callback", "callback", 1)

		// 构造完整跳转 URL
		redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)

		// 执行 302 跳转
		c.Ctx.Redirect(redirectURL, iris.StatusFound)
		return models.OAuthToken{}

	default:
		c.Ctx.StatusCode(iris.StatusBadRequest)
		c.Ctx.WriteString("不支持的 OAuth provider")
		return models.OAuthToken{}
	}
}
