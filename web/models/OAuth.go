package models

// OAuthConfig 用于映射 OAuth 提供商配置
type OAuthConfig struct {
	TeamID      string `json:"team_id"`
	ClientID    string `json:"client_id"`
	KeyID       string `json:"key_id"`
	P8Path      string `json:"p8_path"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
	Oauth_Domain string `json:"oauth_domain"`
}

// OAuthClientConfig 返回给前端的配置
type OAuthClientConfig struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	ClientSecret string `json:"client_secret"`
}

// OAuthToken 返回给前端的配置
type OAuthToken struct {
	Token string `json:"token"`
}

// AppleOAuthUser 对应 Apple 返回的 payload
type AppleOAuthUser struct {
	ISS            string `json:"iss"`              // 签发者
	AUD            string `json:"aud"`              // 接收方 (App 的 bundle ID)
	EXP            int64  `json:"exp"`              // 过期时间
	IAT            int64  `json:"iat"`              // 签发时间
	SUB            string `json:"sub"`              // 用户唯一标识
	ATHash         string `json:"at_hash"`          // Access Token Hash
	Email          string `json:"email"`            // Apple 账号(匿名邮箱)
	EmailVerified  bool   `json:"email_verified"`   // 邮箱是否验证
	IsPrivateEmail bool   `json:"is_private_email"` // 是否为隐藏邮箱
	AuthTime       int64  `json:"auth_time"`        // 认证时间
	NonceSupported bool   `json:"nonce_supported"`  // 是否支持 nonce
}
