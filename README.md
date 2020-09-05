## YF-FPM-SERVER-GO

fpm-server write by golang.

it's a simple framework for implement a go server.

powerful and easy.

## Install

`go get -u github.com/team4yf/yf-fpm-server-go`

## Usage

```golang
package main

import (
	"errors"
	"time"
	"fmt"
	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/pkg/log"
)

func main() {

	app := fpm.New()

	app.Init()

	bizModule := make(fpm.BizModule, 0)
	bizModule["bar"] = func(param *fpm.BizParam) (data interface{}, err error) {
		err = errors.New("foo stub")
		time.Sleep(1 * time.Second)
		data = 1
		return
	}
	app.AddBizModule("foo", &bizModule)

	app.Run()
}

```

now try to execute the biz

`$ curl localhost:9090/biz/foo/bar` 

simple and speed.