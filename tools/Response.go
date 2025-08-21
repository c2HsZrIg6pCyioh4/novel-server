package tools

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// 失败响应
func Fail(code int) Response {
	return Response{
		Code:    code,
		Message: GetMessage(code),
	}
}
