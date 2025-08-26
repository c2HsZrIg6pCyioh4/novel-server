package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
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
type AuthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// POST /oauth/{provider:string}/callback
func (c *Auth_Callback_Controller) Post(provider string) models.OAuthToken {
	var req AuthCallbackRequest
	if provider == "apple" {
		if err := c.Ctx.ReadJSON(&req); err != nil {
			fmt.Println("解析回调参数失败:", err)
			return models.OAuthToken{Token: ""}
		}
		tokenResponse, _ := tools.AppleExchangeCodeForToken(req.Code)
		apple_user, _ := tools.AppleDecodeIDToken(tokenResponse["access_token"].(string))
		use, _ := tools.MySQLGetOpenapiUserbyApplesub(apple_user.SUB)
		tempToken, _ := tools.GenerateJWT(use.Sub, 4) // 4小时有效
		return models.OAuthToken{
			Token: tempToken,
		}
	}

	// 其他 provider 可以自行处理
	return models.OAuthToken{
		Token: "",
	}
}
