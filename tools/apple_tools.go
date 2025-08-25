package tools

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 读取 p8 私钥
func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	var config, _ = GetAppConfig("config.yaml")
	data, err := os.ReadFile(config.OAuth["apple"].P8Path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM format")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not ECDSA")
	}
	return ecdsaKey, nil
}

// GenerateAppleClientSecret 生成 client_secret (默认有效期 4小时)
func GenerateAppleClientSecret(expireSeconds int64) (string, error) {
	var config, _ = GetAppConfig("config.yaml")
	if expireSeconds == 0 {
		expireSeconds = 3600 * 4
	}
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss": config.OAuth["apple"].TeamID,
		"aud": "https://appleid.apple.com",
		"sub": config.OAuth["apple"].ClientID,
		"iat": now,
		"exp": now + expireSeconds,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = config.OAuth["apple"].KeyID
	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return clientSecret, nil
}

// ExchangeCodeForToken 使用 authorization code 换取 token
func ExchangeCodeForToken(code string) (map[string]interface{}, error) {
	var config, _ = GetAppConfig("config.yaml")
	clientSecret, err := GenerateAppleClientSecret(0)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("client_id", config.OAuth["apple"].ClientID)
	values.Set("client_secret", clientSecret)
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", config.OAuth["apple"].RedirectURI)

	resp, err := http.PostForm("https://appleid.apple.com/auth/token", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Apple token request failed: %s", string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return nil, err
	}
	return tokenData, nil
}
