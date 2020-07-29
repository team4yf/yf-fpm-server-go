package main

import (
	"fmt"

	"github.com/team4yf/yf-fpm-server-go/fpm"
	"github.com/team4yf/yf-fpm-server-go/handler"
)

func main() {

	app := &fpm.Fpm{}

	app.New()

	app.AddHook("BEFORE_INIT", func(f *fpm.Fpm) {
		fmt.Println("run some hook")
	}, 10)

	app.Init()

	app.BindHandler("/health", handler.Health).Methods("GET")

	app.BindHandler("/api", handler.API).Methods("POST")

	app.Run(":9999")

}
