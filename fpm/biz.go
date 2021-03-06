package fpm

import (
	"time"

	"github.com/team4yf/fpm-go-pkg/utils"
)

//APIReq api 请求体
type APIReq struct {
	Method    string      `json:"method"`
	Appkey    string      `json:"appkey"`
	Timestamp int64       `json:"timestamp"`
	V         string      `json:"v"`
	Raw       interface{} `json:"param"`
	Sign      string      `json:"sign"`
	Param     *BizParam
}

//APIRsp api 响应体
type APIRsp struct {
	Errno     int         `json:"errno"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

//ResponseOK 构建一个ok的响应体
func ResponseOK(data interface{}) APIRsp {
	return APIRsp{
		Errno:     0,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

//BizParam 业务的请求参数
type BizParam map[string]interface{}

//Convert 将参数转换成实体对象
func (p *BizParam) Convert(obj interface{}) error {
	return utils.Interface2Struct(p, &obj)
}

//BizHandler 具体的业务处理函数
type BizHandler func(*BizParam) (interface{}, error)

//BizModule 业务函数组
type BizModule map[string]BizHandler
