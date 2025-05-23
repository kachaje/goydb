package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/kachaje/goydb/pkg/goydb"
	"github.com/kachaje/goydb/pkg/public"
	"github.com/kachaje/utils"
)

func main() {
	cfg, err := goydb.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.ParseFlags()
	cfg.Containers = []public.Container{
		utils.Fauxton{},
	}

	gdb, err := cfg.BuildDatabase()
	if err != nil {
		log.Fatal(err)
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, gdb.Handler)

	log.Printf("Listening on %s...", cfg.ListenAddress)
	err = http.ListenAndServe(cfg.ListenAddress, loggedRouter)
	if err != nil {
		log.Fatal(err)
	}
}
