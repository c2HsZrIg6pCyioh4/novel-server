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
	config, _ := tools.GetAppConfig("config.yaml")
	if provider == "apple" {
		code := c.Ctx.FormValue("code")
		state := c.Ctx.FormValue("state")
		redirectURI := config.OAuth["apple"].RedirectURI
		redirectURI = strings.Replace(redirectURI, "post_callback", "callback", 1)
		redirectURL := fmt.Sprintf(redirectURI+"?code=%s&state=%s", code, state)
		c.Ctx.Redirect(redirectURL, iris.StatusFound)
		return models.OAuthToken{}
	}
	return models.OAuthToken{}
}
