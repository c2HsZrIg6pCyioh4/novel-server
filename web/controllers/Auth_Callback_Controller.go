package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"novel-server/tools"
	"novel-server/web/models"
	"time"

	"github.com/kataras/iris/v12"
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
			log.Println("解析回调参数失败:", err)
			return models.OAuthToken{Token: ""}
		}
		tokenResponse, _ := tools.AppleExchangeCodeForToken(req.Code)
		idTokenStr, ok := tokenResponse["id_token"].(string)
		if !ok || idTokenStr == "" {
			log.Println("id_token 为空或类型错误:", tokenResponse)
			return models.OAuthToken{Token: ""}
		}
		apple_user, err := tools.AppleDecodeIDToken(idTokenStr)
		if err != nil {
			log.Println("解析 Apple ID Token 失败:", err)
			return models.OAuthToken{Token: ""}
		}
		log.Println(apple_user.SUB)
		// 3. 查询数据库中是否存在该用户
		use, status := tools.MySQLGetOpenapiUserbyApplesub(apple_user.SUB)
		if !status {
			newUser := models.User{
				Sub:      apple_user.SUB,
				Email:    apple_user.Email,
				AppleSub: apple_user.SUB,
			}
			tools.MySQLCreateOpenapiUser(newUser)
			use, status = tools.MySQLGetOpenapiUserbyApplesub(apple_user.SUB)
		}
		tempToken, _ := tools.GenerateJWT(use.Sub, 4) // 4小时有效
		return models.OAuthToken{
			Token: tempToken,
		}
	}
	if provider == "generic" {
		if err := c.Ctx.ReadJSON(&req); err != nil {
			log.Println("解析回调参数失败:", err)
			return models.OAuthToken{Token: ""}
		}
		config, _ := tools.GetAppConfig("config.yaml")
		tokenEndpoint := "https://" + config.OAuth["generic"].Oauth_Domain + "/oauth/token"
		clientID := config.OAuth["generic"].ClientID
		clientSecret := config.OAuth["generic"].TeamID
		redirectURI := config.OAuth["generic"].RedirectURI
		log.Printf("tokenEndpoint: %s, clientID: %s, clientSecret: %s, code: %s, redirectURI: %s", tokenEndpoint, clientID, clientSecret, req.Code, redirectURI)

		tokenResponse, err := ExchangeCodeForToken(tokenEndpoint, clientID, clientSecret, req.Code, redirectURI)
		if err != nil {
			log.Println("获取 token 失败:", err)
			return models.OAuthToken{Token: ""}
		}
		userInfoEndpoint := "https://" + config.OAuth["generic"].Oauth_Domain + "/oauth/userinfo" // 确认具体路径
		user, err := FetchUserInfo(userInfoEndpoint, tokenResponse.AccessToken)
		if err != nil {
			log.Println("获取用户信息失败:", err)
			return models.OAuthToken{Token: ""}
		} else {
			log.Printf("User: %+v", user)
		}
		tempToken, _ := tools.GenerateJWT(user.Sub, 4) // 4小时有效
		return models.OAuthToken{
			Token: tempToken,
		}
	}
	// 其他 provider 可以自行处理
	return models.OAuthToken{
		Token: "",
	}
}

// TokenInfo 令牌信息结构
type TokenInfo struct {
	ClientID    string    `json:"client_id"`
	UserID      string    `json:"user_id"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresIn   int       `json:"expires_in"`
	Scope       string    `json:"scope"`
}

// ExchangeCodeForToken 使用 OAuth2 授权码换取 Token (通用版)
func ExchangeCodeForToken(tokenEndpoint, clientID, clientSecret, code, redirectURI string) (*TokenInfo, error) {
	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("client_secret", clientSecret)
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", redirectURI)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(tokenEndpoint, values)
	if err != nil {
		return nil, fmt.Errorf("请求 token 接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 token 响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token 请求失败, status=%d, body=%s", resp.StatusCode, string(body))
	}

	var tokenInfo TokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("解析 token JSON 失败: %w, body=%s", err, string(body))
	}
	return &tokenInfo, nil
}

// User 用户模型
type User struct {
	UnionID       string `json:"unionid"`
	Nickname      string `json:"nickname"`
	PhoneNumber   string `json:"phonenumber"`
	Email         string `json:"email"`
	RandNum       string `json:"randnum"`
	EnabledStatus int    `json:"enabled_status"`
	IsActive      bool   `json:"is_active"`
	Username      string `json:"username"`
	Sub           string `json:"sub"`
	AppleSub      string `json:"apple_sub"`
	Password      string `json:"password"`
}

// FetchUserInfo 使用 AccessToken 获取用户信息
func FetchUserInfo(userInfoEndpoint, accessToken string) (*User, error) {
	log.Println("获取用户信息accessToken:", accessToken)
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("构建请求失败: %w", err)
	}

	// 设置 Bearer Token
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求用户信息接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取用户信息响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("用户信息请求失败, status=%d, body=%s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("解析用户信息 JSON 失败: %w, body=%s", err, string(body))
	}

	return &user, nil
}
