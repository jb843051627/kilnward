package handler

import (
	_ "embed"
	"net/http"
)

//go:embed assets/index.html
var indexHTML []byte

func serveIndex(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}

func init() { _ = http.MethodGet }
