package main

import "go.uber.org/fx/service"

func main() {
	svc, _ := service.WithModule(myFunc).Build()
	svc.Start()
}
