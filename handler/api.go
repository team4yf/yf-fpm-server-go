//Package handler the default handlers
package handler

import "github.com/team4yf/yf-fpm-server-go/core"

func Health(ctx *core.Ctx) {

	ctx.JSON(map[string]interface{}{"Status": "UP"})
}

func API(ctx *core.Ctx) {
	// an example API handler
	data := make(map[string]interface{}, 0)
	if err := ctx.ParseBody(&data); err != nil {
		return
	}

	ctx.JSON(data)
}
