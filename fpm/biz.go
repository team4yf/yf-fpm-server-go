package fpm

//APIReq api 请求体
type APIReq struct {
	Method    string    `json:"method"`
	Appkey    string    `json:"appkey"`
	Timestamp int64     `json:"timestamp"`
	V         string    `json:"v"`
	Param     *BizParam `json:"param"`
	Sign      string    `json:"sign"`
}

//APIRsp api 响应体
type APIRsp struct {
	Errno     int         `json:"errno"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

//BizParam 业务的请求参数
type BizParam map[string]interface{}

//BizHandler 具体的业务处理函数
type BizHandler func(*BizParam) (interface{}, error)

//BizModule 业务函数组
type BizModule map[string]BizHandler
