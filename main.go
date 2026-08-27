package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jb843051627/kilnward/internal/handler"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func main() {
	dbPath := os.Getenv("KILNWARD_DB")
	if dbPath == "" {
		dbPath = "data/kilnward.db"
	}
	repository, err := store.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()

	app := service.NewLab(repository)
	defer app.Close()

	addr := os.Getenv("KILNWARD_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("kilnward listening on %s", addr)
	server := &http.Server{Addr: addr, Handler: handler.New(app)}
	log.Fatal(server.ListenAndServe())
}
