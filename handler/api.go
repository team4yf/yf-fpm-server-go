//Package handler the default handlers
package handler

import "github.com/team4yf/yf-fpm-server-go/fpm"

func Health(ctx *fpm.Ctx) {

	ctx.JSON(map[string]interface{}{"Status": "UP"})
}

func API(ctx *fpm.Ctx) {
	// an example API handler
	data := make(map[string]interface{}, 0)
	if err := ctx.ParseBody(&data); err != nil {
		return
	}

	ctx.JSON(data)
}
