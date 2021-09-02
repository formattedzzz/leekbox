package model

type Resp struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	API_SUCCESS     = "请求成功"
	UNHANDLED_ERROR = "发生未知错误"
)
