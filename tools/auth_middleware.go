package tools

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

// TokenAuthMiddleware 返回一个支持全局配置的 Token 中间件
// enableHealthAuth: 是否对 /health 接口也进行认证
// skipGET: 是否跳过 GET 请求
func TokenAuthMiddleware(validTokens []string, enableAuth bool) iris.Handler {
	// 使用 map 提升查找效率
	tokenMap := make(map[string]struct{}, len(validTokens))
	for _, t := range validTokens {
		tokenMap[t] = struct{}{}
	}

	return func(ctx iris.Context) {
		if !enableAuth {
			ctx.Next()
			return
		}

		if ctx.Path() == "/health" ||
			(strings.HasPrefix(ctx.Path(), "/oauth/") &&
				(strings.HasSuffix(ctx.Path(), "/callback") || strings.HasSuffix(ctx.Path(), "/post_callback"))) ||
			ctx.Method() == http.MethodGet {
			ctx.Next()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			sendFail(ctx, ErrorCode.AuthHeaderMissing)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			sendFail(ctx, ErrorCode.AuthTokenFormatWrong)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// 1. 先验证本地静态 token
		if _, ok := tokenMap[tokenString]; ok {
			ctx.Next()
			return
		}

		// 2. 再验证 JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			config, _ := GetAppConfig("config.yaml")
			return []byte(config.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			sendFail(ctx, ErrorCode.AuthTokenInvalid)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Values().Set("user_sub", claims["sub"])
		}

		ctx.Next()
	}
}

// sendFail 统一返回 Response
func sendFail(ctx iris.Context, code int) {
	ctx.StatusCode(http.StatusUnauthorized)
	ctx.JSON(Fail(code)) // 使用全局 Fail 函数，保证统一格式
	ctx.StopExecution()
}

// GenerateJWT 生成短期 JWT token
func GenerateJWT(userID string, expireHours int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // HS256
	config, _ := GetAppConfig("config.yaml")
	return token.SignedString([]byte(config.JwtSecret)) // 用字符串即可
}
