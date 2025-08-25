package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"log"
	"novel-server/tools"
	"novel-server/web/models"
)

type Auth_Callback_Controller struct {
	Ctx iris.Context
}

type AppleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// POST /oauth/{provider:string}/callback
func (c *Auth_Callback_Controller) Post(provider string) models.OAuthToken {
	if provider == "apple" {
		// 1. 读取回调参数
		code := c.Ctx.FormValue("code")
		state := c.Ctx.FormValue("state")
		userStr := c.Ctx.FormValue("user")
		fmt.Println("state:", state)
		fmt.Println("code:", code)
		fmt.Println("user raw:", userStr)

		// 2. TODO: 用 code 调 Apple 接口换取 token
		tokenResponse, _ := tools.ExchangeCodeForToken(code)
		// 这里先演示返回 id_token
		log.Print(tokenResponse)
		idToken := "sk-123456"
		return models.OAuthToken{
			Token: idToken,
		}
	}

	// 其他 provider 可以自行处理
	return models.OAuthToken{
		Token: "",
	}
}
