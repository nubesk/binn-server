package main

import (
	"log"

	"github.com/nubesk/binn"
	"github.com/nubesk/binn-server/config"
	"github.com/nubesk/binn-server/server"
)

func main() {
	c := config.NewFromEnv()
	bn := binn.New(binn.NewBottleStorage(100), c.SendInterval)
	l := log.Default()
	srv := server.New(bn, ":8080", l)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
