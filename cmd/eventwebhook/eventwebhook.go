package main

import (
	"log"

	"github.com/jimmyjames85/SendGrid/pkg/server"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	cfg := server.Config{}
	envconfig.MustProcess("EVENTWEBHOOK", &cfg)
	log.Printf("%s", cfg.ToJSON())
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("err loading config: %v", err)
	}
	err = srv.Serve()
	if err != nil {
		log.Fatalf("server failure: %v", err)
	}
}
