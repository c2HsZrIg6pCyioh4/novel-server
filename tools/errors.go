package tools

// ErrorCode 用作命名空间管理错误码
var ErrorCode = struct {
	CodeSuccess            int
	CodeDeleteNovelFailed  int
	CodeGetLatestIndexFail int
	CodeNotLatestChapter   int
	CodeDeleteChapterFail  int
	AuthHeaderMissing    int
	AuthTokenInvalid     int
	AuthTokenFormatWrong int
}{
	CodeSuccess:            0,
	CodeDeleteNovelFailed:  1001,
	CodeGetLatestIndexFail: 1002,
	CodeNotLatestChapter:   1003,
	CodeDeleteChapterFail:  1004,
	AuthHeaderMissing:    2001,
	AuthTokenInvalid:     2002,
	AuthTokenFormatWrong: 2003,
}

// 错误码对应信息
var codeMessages = map[int]string{
	0:    "success",
	1001: "书籍删除失败",
	1002: "获取最新章节失败",
	1003: "只能删除最新章节",
	1004: "删除章节失败",
	ErrorCode.AuthHeaderMissing:    "Authorization header missing",
	ErrorCode.AuthTokenInvalid:     "Invalid token",
	ErrorCode.AuthTokenFormatWrong: "Invalid token format",
}

// GetMessage 根据错误码获取错误信息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
