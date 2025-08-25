package controllers

import (
	"github.com/kataras/iris/v12"
	"novel-server/tools"
	"novel-server/web/models"
)

type Auth_Config_Controller struct {
	Ctx iris.Context
}

// Get "/oauth/{provider:string}/config"
func (c *Auth_Config_Controller) Get(provider string) models.OAuthClientConfig {
	config, _ := tools.GetAppConfig("config.yaml")

	data, _ := config.OAuth[provider]

	clientConfig := models.OAuthClientConfig{
		ClientID:    data.ClientID,
		RedirectURI: data.RedirectURI,
		Scope:       data.Scope,
	}
	return clientConfig
}
