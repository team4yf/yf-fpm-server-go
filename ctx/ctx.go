package ctx

import (
	"encoding/json"
	"net/http"

	"github.com/team4yf/yf-fpm-server-go/pkg/utils"
)

//Ctx the content of the request
type Ctx struct {
	request *http.Request
	w       http.ResponseWriter
}

//WrapCtx wrap the context with the w & request
func WrapCtx(w http.ResponseWriter, request *http.Request) *Ctx {

	return &Ctx{
		request: request,
		w:       w,
	}
}

//JSON output the json
func (c *Ctx) JSON(data interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(http.StatusOK)
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

//Fail execute fail
func (c *Ctx) Fail(err interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(http.StatusOK)
	json.NewEncoder(c.w).Encode(err)
}

//GetRequest output the json
func (c *Ctx) GetRequest() *http.Request {
	return c.request
}
