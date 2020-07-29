package main

import (
	"fmt"

	"github.com/team4yf/yf-fpm-server-go/core"
	"github.com/team4yf/yf-fpm-server-go/handler"
)

func main() {

	fpm := &core.Fpm{}

	fpm.New()

	fpm.AddHook("BEFORE_INIT", func(f *core.Fpm) {
		fmt.Println("run some hook")
	}, 10)

	fpm.Init()

	fpm.BindHandler("/health", handler.Health).Methods("GET")

	fpm.BindHandler("/api", handler.API).Methods("POST")

	fpm.Run(":9999")

}
