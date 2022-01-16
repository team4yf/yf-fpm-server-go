package fpm

import (
	"time"

	"github.com/team4yf/fpm-go-pkg/utils"
)

//APIReq api request body
type APIReq struct {
	Method    string      `json:"method"`
	Appkey    string      `json:"appkey"`
	Timestamp int64       `json:"timestamp"`
	V         string      `json:"v"`
	Raw       interface{} `json:"param"`
	Sign      string      `json:"sign"`
	Param     *BizParam
}

//APIRsp api response body
type APIRsp struct {
	Errno     int         `json:"errno"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

//ResponseOK create a success response
func ResponseOK(data interface{}) APIRsp {
	return APIRsp{
		Errno:     0,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

type IBizParam interface {
	Pre() map[int]interface{}
	GetResult() interface{}
	SetResult(interface{})
	Post() map[int]interface{}
	Convert(interface{}) error
}

//BizParam method param
type BizParam map[string]interface{}

// type BizParam struct {
// 	__pre__    map[int]interface{}
// 	__result__ interface{} `json:"__result__,omitempty"`
// 	__post__   map[int]interface{}
// }

//Convert 将参数转换成实体对象
func (p *BizParam) Convert(obj interface{}) error {
	return utils.Interface2Struct(p, &obj)
}

func (p *BizParam) GetResult() interface{} {
	return (*p)["__result__"]
}

func (p *BizParam) SetResult(obj interface{}) {
	(*p)["__result__"] = obj
}

func (p *BizParam) Pre() map[int]interface{} {
	if pre, ok := (*p)["__pre__"]; ok {
		return pre.(map[int]interface{})
	}
	pre := make(map[int]interface{})
	(*p)["__pre__"] = pre
	return pre
}

func (p *BizParam) Post() map[int]interface{} {
	if post, ok := (*p)["__post__"]; ok {
		return post.(map[int]interface{})
	}
	post := make(map[int]interface{})
	(*p)["__post__"] = post
	return post
}

//BizHandler 具体的业务处理函数
type BizHandler func(*BizParam) (interface{}, error)

//BizModule 业务函数组
type BizModule map[string]BizHandler
