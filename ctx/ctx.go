package ctx

import (
	"encoding/json"
	"net/http"

	"strings"

	"github.com/gorilla/mux"
	"github.com/team4yf/fpm-go-pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/errno"
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

//Param get the url path from the request
func (c *Ctx) Param(p string) string {
	vars := mux.Vars(c.request)
	return vars[p]
}

//Query get the query string of the url
func (c *Ctx) Query(p string) string {
	q := c.request.URL.Query()
	return q.Get(p)

}

//Querys get the query string of the url
func (c *Ctx) Querys() map[string]string {
	querys := make(map[string]string)
	for k, vs := range c.request.URL.Query() {
		if len(vs) < 1 {
			querys[k] = ""
		} else {
			querys[k] = vs[0]
		}

	}
	return querys

}

//GetHeader get data from the header
func (c *Ctx) GetHeader(k string) string {
	data := c.request.Header.Get(k)
	return data
}

//GetRemoteIP get remote ip from request
func (c *Ctx) GetRemoteIP() string {
	for _, k := range []string{"X-Real-Ip", "X-Real-IP", "X-FORWARDED-FOR", "X-Forwarded-For", "x-forwarded-for"} {
		if id := c.request.Header.Get(k); id != "" {
			return id
		}
	}
	return c.request.RemoteAddr
}

//GetRequestID get request id from request header
func (c *Ctx) GetRequestID() (id string) {
	// try to get request-id from
	for _, k := range []string{"HTTP_X_REQUEST_ID", "X-Request-ID", "X-Request-Id", "x-request-id"} {
		if id = c.request.Header.Get(k); id != "" {
			return
		}
	}
	return ""
}

//GetToken get token from request header
// it's `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzEyMzEyMyJ9.5KEyBxDH2NGVKoiA0J6IPB4QPlvZi9zPH9SSKTWF2h8` in header[`Authorization`]
func (c *Ctx) GetToken() string {
	token := c.request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		return token[7:]
	}
	return token
}

//QueryDefault get the querystring of the url, return default value if nil
func (c *Ctx) QueryDefault(p, dfv string) string {
	q := c.request.URL.Query()
	data := q.Get(p)
	if data == "" {
		return dfv
	}
	return data
}

//JSON output the json
func (c *Ctx) JSON(data interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(http.StatusOK)
	json.NewEncoder(c.w).Encode(data)
}

//JSONWithoutHTMLEscape output the json without html escaping
func (c *Ctx) JSONWithoutHTMLEscape(data interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(c.w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
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

//GetResponse output the json
func (c *Ctx) GetResponse() http.ResponseWriter {
	return c.w
}

//BizError output the biz error
func (c *Ctx) BizError(err *errno.BizError) {
	c.JSON(err)
}
