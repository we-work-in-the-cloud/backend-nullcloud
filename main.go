package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func main() {
	storeFile := flag.String("store-file", "", "path to JSON persistence file (default: in-memory)")
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	var s store.Store
	if *storeFile != "" {
		fs, err := store.NewJSONFileStore(*storeFile)
		if err != nil {
			log.Fatalf("failed to open store file: %v", err)
		}
		log.Printf("using JSON file store: %s", *storeFile)
		s = fs
	} else {
		s = store.NewMemoryStore()
		log.Println("using in-memory store")
	}

	srv := api.NewServer(s)
	addr := ":" + *port
	log.Printf("starting NullCloud backend on %s", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
