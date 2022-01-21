package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/team4yf/fpm-go-pkg/log"
	"github.com/team4yf/yf-fpm-server-go/ctx"
)

type AspectLogConfig struct {
	Enable  bool
	App     string
	Pattern []string
}

type bodyLogWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AspectLog(aspectLogConfig *AspectLogConfig) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !aspectLogConfig.Enable || !matchURL(r.URL.String(), aspectLogConfig.Pattern) {
				next.ServeHTTP(w, r)
				return
			}
			bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: w}
			w = bodyLogWriter

			var requestData string
			contentType := r.Header.Get("Content-Type")
			if r.Method == "POST" {
				if strings.Contains(contentType, "json") {

					if r.Body != nil {
						data, _ := ioutil.ReadAll(r.Body)
						// 这里需要将原来的数据还原回去，否则后面的handler获取不到原来的请求数据
						r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
						requestData = string(data)
					}
				}
			}

			//开始时间
			startTime := time.Now()

			//处理请求
			next.ServeHTTP(w, r)

			responseBody := bodyLogWriter.body.String()

			//结束时间
			endTime := time.Now()

			//日志格式
			accessLogMap := make(map[string]interface{})

			accessLogMap["app"] = aspectLogConfig.App
			accessLogMap["req_time"] = startTime
			accessLogMap["req_method"] = r.Method
			accessLogMap["req_uri"] = r.RequestURI
			accessLogMap["req_ua"] = r.UserAgent()
			accessLogMap["req_content_type"] = contentType
			accessLogMap["req_referer"] = r.Referer()
			accessLogMap["req_post_data"] = fmt.Sprintf("Body: %s", requestData)
			// try to get client IP
			accessLogMap["req_client_ip"] = ctx.WrapCtx(w, r).GetRemoteIP()

			accessLogMap["rsp_time"] = endTime
			accessLogMap["rsp_body"] = responseBody

			accessLogMap["cost_time"] = fmt.Sprintf("%v", endTime.Sub(startTime))

			accessLogJSON, _ := json.Marshal(accessLogMap)

			log.Info(string(accessLogJSON))
		}

		return http.HandlerFunc(fn)
	}
}
