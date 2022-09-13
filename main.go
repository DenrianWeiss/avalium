package main

import (
	"github.com/DenrianWeiss/avalium/config"
	"github.com/DenrianWeiss/avalium/model"
	"github.com/DenrianWeiss/avalium/service/forward"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	// Initialize Config
	config.Init()
	model.GenerateCache(config.GetConfig().Rpc)
	log.Println("Listening on", config.GetConfig().Serve.ServerAddr)
	err := fasthttp.ListenAndServe(config.GetConfig().Serve.ServerAddr, forward.ReceiveHttp)
	if err != nil {
		panic(err)
	}
}
