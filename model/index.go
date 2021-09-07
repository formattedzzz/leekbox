package model

type Resp struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Return(code int, data interface{}, message string) Resp {
	return Resp{code, data, message}
}

const (
	API_SUCCESS     = "请求成功"
	UNHANDLED_ERROR = "发生未知错误"
	PARAMS_ERROR    = "请检查请求参数"
)
