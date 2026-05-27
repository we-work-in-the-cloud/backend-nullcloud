package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

var version = "dev"

func main() {
	storeFile := flag.String("store-file", "", "path to JSON persistence file (default: in-memory)")
	port := flag.String("port", "8080", "port to listen on")
	tokensFlag := flag.String("tokens", "", "comma-separated list of allowed API tokens (default: all tokens allowed)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	var allowedTokens []string
	if *tokensFlag != "" {
		for _, t := range strings.Split(*tokensFlag, ",") {
			if t = strings.TrimSpace(t); t != "" {
				allowedTokens = append(allowedTokens, t)
			}
		}
	}

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

	srv := api.NewServer(s, allowedTokens)
	addr := ":" + *port
	log.Printf("starting NullCloud backend on %s", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
