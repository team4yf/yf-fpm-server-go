package core

import (
	"encoding/json"
	"net/http"

	"github.com/team4yf/yf-fpm-server-go/pkg/utils"
)

//Ctx the content of the request
type Ctx struct {
	fpm     *Fpm
	request *http.Request
	w       http.ResponseWriter
}

//WrapCtx wrap the context with the w & request
func WrapCtx(fpm *Fpm, w http.ResponseWriter, request *http.Request) *Ctx {

	return &Ctx{
		fpm:     fpm,
		request: request,
		w:       w,
	}
}

//JSON output the json
func (c *Ctx) JSON(data interface{}) {
	json.NewEncoder(c.w).Encode(data)
}

//ParseBody parse the request body to the data
func (c *Ctx) ParseBody(data interface{}) (err error) {
	err = utils.GetBodyStruct(c.request.Body, data)

	if err != nil {
		return
	}
	return
}

//GetRequest output the json
func (c *Ctx) GetRequest() *http.Request {
	return c.request
}
