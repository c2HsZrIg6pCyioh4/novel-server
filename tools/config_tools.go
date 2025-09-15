package tools

import (
	"encoding/json"
	"novel-server/web/models"
	"os"
)

// AppConfig 结构体用于映射JSON配置文件
type AppConfig struct {
	AppName    string `json:"app_name"`
	Port       string `json:"port"`
	StaticPath string `json:"static_path"`
	Redis      struct {
		Network  string `json:"network"`
		Addr     string `json:"addr"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Prefix   string `json:"prefix"`
	} `json:"redis"`
	MySQL struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	} `json:"mysql"`
	OpenApiMySQL struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	} `json:"openapi_mysql"`
	Weixin struct {
		URL            string `json:"url"`
		AppID          string `json:"app_id"`
		AppSecret      string `json:"app_secret"`
		Token          string `json:"token"`
		EncodingAESKey string `json:"encodingaeskey"`
		Env_Version    string `json:"env_version"`
	} `json:"weixin"`
	Openai struct {
		URL            string `json:"url"`
		AuthToken      string `json:"authtoken"`
		FREE_URL       string `json:"free_url"`
		FREE_AuthToken string `json:"free_authtoken"`
	} `json:"openai"`
	GeoLite struct {
		CityDBPath    string `json:"citydbpath"`
		CountryDBPath string `json:"countrydbpath"`
		ASNDBPath     string `json:"asndbpath"`
	} `json:"geolite"`
	Mode        string   `json:"mode"`
	LoggerLevel string   `json:"loggerlevel"`
	Tokens      []string `json:"Tokens"` // 修正这里
	// Token 中间件相关配置
	Auth_status bool `json:"auth_Status"` // 是否启用 Token 认证
	// ---------- 新增 OAuth 配置 ----------
	OAuth     map[string]models.OAuthConfig `json:"oauth"` // key: "apple", "google"
	JwtSecret string                        `json:"jwtsecret"`
}

var AppConfigInstance *AppConfig

// GetAppConfig 从JSON文件中获取配置信息
func GetAppConfig(filePath string) (*AppConfig, error) {
	configData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析JSON配置
	var config AppConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	// 设置全局 AppConfigInstance
	AppConfigInstance = &config
	return &config, nil
}
