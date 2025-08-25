package models

// OAuthConfig 用于映射 OAuth 提供商配置
type OAuthConfig struct {
	TeamID      string `json:"team_id"`
	ClientID    string `json:"client_id"`
	KeyID       string `json:"key_id"`
	P8Path      string `json:"p8_path"`
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
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
