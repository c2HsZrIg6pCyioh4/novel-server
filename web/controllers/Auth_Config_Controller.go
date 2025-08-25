package controllers

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"novel-server/tools"
	"novel-server/web/models"
	"os"
	"time"

	"github.com/kataras/iris/v12"
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
	// Apple 需要生成 client_secret
	if provider == "apple" {
		clientSecret, _ := GenerateAppleClientSecret(data.P8Path, data.TeamID, data.ClientID, data.KeyID)
		clientConfig.ClientSecret = clientSecret
	}
	return clientConfig
}

// GenerateAppleClientSecret 生成 Apple client_secret
func GenerateAppleClientSecret(p8Path, teamID, clientID, keyID string) (string, error) {
	p8Data, err := os.ReadFile(p8Path)
	if err != nil {
		return "", fmt.Errorf("read p8 file failed: %w", err)
	}

	block, _ := pem.Decode(p8Data)
	if block == nil {
		return "", fmt.Errorf("invalid PEM format")
	}

	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse private key failed: %w", err)
	}

	privateKey, ok := privateKeyInterface.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not ECDSA")
	}

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss": teamID,
		"iat": now,
		"exp": now + 86400*180, // 最大 6 个月
		"aud": "https://appleid.apple.com",
		"sub": clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = keyID

	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign JWT failed: %w", err)
	}

	return clientSecret, nil
}
