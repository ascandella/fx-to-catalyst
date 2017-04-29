// taken from github.com/uber/go-fx/examples/simple/main.go

package main

import (
	"go.uber.org/fx/modules/uhttp"
	"go.uber.org/fx/service"
)

func main() {
	svc, _ := service.WithModule(
		uhttp.New(nil),
	).Build()

	svc.Start()
}
