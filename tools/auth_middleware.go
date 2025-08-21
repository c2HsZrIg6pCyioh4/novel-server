package tools

import (
	"net/http"
	"strings"

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
		// 路径和方法判断
		if ctx.Path() == "/health" {
			ctx.Next()
			return
		}

		if ctx.Method() == http.MethodGet {
			ctx.Next()
			return
		}

		// 获取 Header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			sendFail(ctx, ErrorCode.AuthHeaderMissing)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			sendFail(ctx, ErrorCode.AuthTokenFormatWrong)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// 校验 Token
		if _, ok := tokenMap[token]; !ok {
			sendFail(ctx, ErrorCode.AuthTokenInvalid)
			return
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
