package api

import (
	_ "embed"
	"net/http"
)

//go:embed ui.html
var uiHTML []byte

//go:embed ui.css
var uiCSS []byte

//go:embed ui.js
var uiJS []byte

func uiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(uiHTML)
	}
}

func uiCSSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(uiCSS)
	}
}

func uiJSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(uiJS)
	}
}
