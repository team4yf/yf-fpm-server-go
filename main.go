package main

import (
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/log"

	_ "github.com/team4yf/yf-fpm-server-go/plugin"
)

func main() {

	app := fpm.New()

	// app.AddHook("BEFORE_INIT", func(f *fpm.Fpm) {
	// 	fmt.Println("run some hook")
	// }, 10)

	app.Init()

	if app.HasConfig("db") {
		dbConfig := app.GetConfig("db")
		log.Debugf("dbconfig %+v", dbConfig)
	}
	bizModule := make(fpm.BizModule, 0)
	bizModule["bar"] = func(param *fpm.BizParam) (data interface{}, err error) {
		data = 1
		return
	}
	app.AddBizModule("foo", &bizModule)

	app.Run(":9999")

}
