package main

import (
	"time"

	"github.com/pkg/errors"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/log"
	_ "github.com/team4yf/yf-fpm-server-go/plugin"
)

type DBSetting struct {
	Engine   string
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Charset  string
	ShowSQL  bool
}

type testBody struct {
	Foo  string
	Shit string
}

func main() {

	app := fpm.New()

	// app.AddHook("BEFORE_INIT", func(f *fpm.Fpm) {
	// 	fmt.Println("run some hook")
	// }, 10)

	app.Init()

	p := &fpm.BizParam{
		"foo":  "bar",
		"shit": "damn",
	}

	body := testBody{}
	err := p.Convert(&body)
	log.Debugf("err: %v, body: %#v", err, body.Foo)
	if app.HasConfig("db") {
		var dbConfig DBSetting
		app.FetchConfig("db", &dbConfig)

		log.Debugf("dbconfig %+v", dbConfig)
	}

	info := app.GetAppInfo()
	log.Debugf("appInfo: %+v", info)
	app.AddFilter("foo.bar", "before", func(app *fpm.Fpm, biz string, args *fpm.BizParam) (bool, error) {
		log.Debugf("before %s, args: %v", biz, *args)
		return true, nil
	}, 1)

	app.AddFilter("foo.bar", "after", func(app *fpm.Fpm, biz string, args *fpm.BizParam) (bool, error) {
		log.Debugf("after %s, args: %v", biz, *args)
		return true, nil
	}, 1)

	app.Subscribe("#webhook/test/foo", func(_ string, data interface{}) {
		log.Debugf("webhook: %+v", data)
	})

	bizModule := make(fpm.BizModule, 0)
	bizModule["bar"] = func(param *fpm.BizParam) (data interface{}, err error) {
		err = errors.New("foo stub")
		time.Sleep(10 * time.Second)
		data = 1
		return
	}
	app.AddBizModule("foo", &bizModule)

	app.Run()

}
