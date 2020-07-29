package main

import (
	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/yf-fpm-server-go/plugin"
)

func main() {

	app := fpm.New()

	// app.AddHook("BEFORE_INIT", func(f *fpm.Fpm) {
	// 	fmt.Println("run some hook")
	// }, 10)

	app.Init()

	bizModule := make(fpm.BizModule, 0)
	bizModule["bar"] = func(param *fpm.BizParam) (data interface{}, err error) {
		data = 1
		return
	}
	app.AddBizModule("foo", &bizModule)

	app.Run(":9999")

}
