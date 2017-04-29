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
	checkError(err, "Unable to initialize service")

	svc.Start()
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, "error", err)
	}
}
