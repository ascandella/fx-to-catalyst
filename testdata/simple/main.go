// taken from github.com/uber/go-fx/examples/simple/main.go

package main

import (
	"log"

	"go.uber.org/fx/modules/uhttp"
	"go.uber.org/fx/service"
)

func main() {
	svc, err := service.WithModule(
		uhttp.New(registerHTTPers, uhttp.WithInboundMiddleware(simpleInboundMiddleware{})),
	).Build()

	if err != nil {
		log.Fatal("Unable to initialize service", "error", err)
	}

	svc.Start()
}
