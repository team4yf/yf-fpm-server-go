package plugin

import (
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func init() {
	fpm.Register(func(app *fpm.Fpm) {
		app.BindHandler("/plugin", func(ctx *ctx.Ctx, _ *fpm.Fpm) {
			ctx.JSON(`{"plugin":"foo1"}`)
		}).Methods("GET")
	})
}
