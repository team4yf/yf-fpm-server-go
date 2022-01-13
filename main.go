package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/team4yf/fpm-go-pkg/log"
	"github.com/team4yf/yf-fpm-server-go/errno"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func main() {

	fpm.RegisterByPlugin(&fpm.Plugin{
		Name: "fpm-plugin-1",
		Handler: func(*fpm.Fpm) {
			fmt.Println("load fpm-plugin-1 ok")
		},
		Deps: []string{"fpm-plugin-2"},
	})

	fpm.RegisterByPlugin(&fpm.Plugin{
		Name: "fpm-plugin-2",
		Handler: func(*fpm.Fpm) {
			fmt.Println("load fpm-plugin-2 ok")
		},
	})

	app := fpm.New()
	app.HealthCheckData.Version = "1"
	app.Init()

	info := app.GetAppInfo()
	log.Debugf("appInfo: %+v", info)

	bizModule := make(fpm.BizModule, 0)
	bizModule["bar"] = func(param *fpm.BizParam) (data interface{}, err error) {
		err = errno.New(-11, "foo stub")
		time.Sleep(1 * time.Second)
		data = 1
		return
	}
	bizModule["echo"] = func(param *fpm.BizParam) (data interface{}, err error) {
		data = "boom!!!"
		return
	}
	app.AddBizModule("foo", &bizModule)

	app.AddFilter("user.login", "before", func(app *fpm.Fpm, biz string, args *fpm.BizParam) (bool, interface{}, error) {
		return true, "userId", errors.New("login failed")
	}, 1)

	app.Run()

}
