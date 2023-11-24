package utils

import (
	"log"
	"net/http"
	"net/http/httptest"
)

func StartLocalServer(path string) *httptest.Server {
	fileServerHandler := http.FileServer(http.Dir(path))
	log.Println("Set file server path:", path)

	server := httptest.NewServer(fileServerHandler)
	log.Println("Server started at:", server.URL)

	return server
}
