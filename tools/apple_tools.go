package tools

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"novel-server/web/models"
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
func AppleExchangeCodeForToken(code string) (map[string]interface{}, error) {
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

	// 设置超时客户端
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm("https://appleid.apple.com/auth/token", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取 Apple token 响应失败: %w", err)
		return nil, nil
	}
	log.Println(resp.StatusCode)
	log.Println(resp.Body)
	if resp.StatusCode != 200 {
		log.Println("Apple token 请求失败, status=%d, body=%s", resp.StatusCode, string(body))
		return nil, nil
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		log.Println("解析 Apple token JSON 失败: %w, body=%s", err, string(body))
		return nil, nil
	}

	// 打印返回结果，便于调试
	log.Printf("Apple token response: %+v\n", tokenData)

	return tokenData, nil
}

// AppleDecodeIDToken 解析 Apple 返回的 id_token，并映射到 OAuthUser
func AppleDecodeIDToken(idToken string) (*models.AppleOAuthUser, error) {
	token, _, err := jwt.NewParser().ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 转成 JSON 再反序列化到结构体
		data, _ := json.Marshal(claims)
		var user models.AppleOAuthUser
		if err := json.Unmarshal(data, &user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, fmt.Errorf("cannot parse claims")
}
