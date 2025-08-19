package tools

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

var sessionManager *sessions.Sessions

// InitSessionManager initializes the session manager.
func InitSessionManager() {
	sessionManager = sessions.New(sessions.Config{
		Cookie:       "openapi_session_id",
		Expires:      3 * time.Minute,
		AllowReclaim: true,
	})
}

// SessionClient returns the session manager.
func SessionClient() *sessions.Sessions {
	if sessionManager == nil {
		InitSessionManager()
	}
	return sessionManager
}

// GetSessionValue retrieves the value for the specified key from session.
//
//	tools.GetSessionValue(ctx, "username")
func GetSessionValue(ctx iris.Context, key string) string {
	session := SessionClient().Start(ctx)
	return session.GetString(key)
}

// SetSessionValue sets the value for the specified key in session.
//
//	tools.SetSessionValue(ctx, "username","111")
func SetSessionValue(ctx iris.Context, key, value string) bool {
	session := SessionClient().Start(ctx)
	if session.Get(key) == nil {
		session.Set(key, value)
		return true
	} else if session.Get(key) == value {
		return true
	} else {
		return false
	}
}

// CheckSessionValue checks if the session value matches the expected value.
//
//	if tools.CheckSessionValue(ctx, "username", "xxxxxx") {
//		println("Session value is valid.")
//	} else {
//		println("Session value is not valid.")
//	}
func CheckSessionValue(ctx iris.Context, key, expectedValue string) bool {
	// 检查 Redis 中是否存在该会话 ID
	if GetSessionIfExists(ctx) != nil {
		session := SessionClient().Start(ctx)
		// 检查会话(session)中是否存在指定键以及其值是否匹配
		if session.GetString(key) == expectedValue {
			return true
		}
		return false
	} else {
		return false
	}

}

// DeleteSession deletes the session.
//
//	tools.DeleteSession(ctx)
func DeleteSession(ctx iris.Context) bool {
	session := SessionClient().Start(ctx)
	// 检查会话(session)中是否存在指定键以及其值是否匹配
	session.Destroy()
	return true
}

func CreateSession(ctx iris.Context) *sessions.Session {
	// 使用会话客户端来创建一个新的会话
	session := sessionManager.Start(ctx)
	return session
}

func GetSessionID(ctx iris.Context) string {
	session := SessionClient().Start(ctx)
	return session.ID()
}

// getSession returns the existing session if it exists, otherwise nil.
func GetSessionIfExists(ctx iris.Context) *sessions.Session {
	if Redis_ExistsValue("openapi_session_id-"+ctx.GetCookie("openapi_session_id")) == 1 {
		session := SessionClient().Start(ctx)
		return session
	}
	return nil
}
