package tools

// 用一个 struct 作为命名空间
var ErrorCode = struct {
	CodeSuccess            int
	CodeDeleteNovelFailed  int
	CodeGetLatestIndexFail int
	CodeNotLatestChapter   int
	CodeDeleteChapterFail  int
}{
	CodeSuccess:            0,
	CodeDeleteNovelFailed:  1001,
	CodeGetLatestIndexFail: 1002,
	CodeNotLatestChapter:   1003,
	CodeDeleteChapterFail:  1004,
}

// 错误码对应信息
var codeMessages = map[int]string{
	0:    "success",
	1001: "书籍删除失败",
	1002: "获取最新章节失败",
	1003: "只能删除最新章节",
	1004: "删除章节失败",
}

// GetMessage 根据错误码获取错误信息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
