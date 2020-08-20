package main

import (
	"time"

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

	app.AddFilter("foo.bar", "before", func(app *fpm.Fpm, biz string, args *fpm.BizParam) (bool, error) {
		log.Debugf("before %s, args: %v", biz, *args)
		return true, nil
	}, 1)

	app.AddFilter("foo.bar", "after", func(app *fpm.Fpm, biz string, args *fpm.BizParam) (bool, error) {
		log.Debugf("after %s, args: %v", biz, *args)
		return true, nil
	}, 1)

	bizModule := make(fpm.BizModule, 0)
	bizModule["bar"] = func(param *fpm.BizParam) (data interface{}, err error) {
		time.Sleep(10 * time.Second)
		data = 1
		return
	}
	app.AddBizModule("foo", &bizModule)

	app.Run()

}
