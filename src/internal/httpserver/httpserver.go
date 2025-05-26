package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

func New(service QuoteService, router *mux.Router, listenAddr string) *http.Server {
	server := &http.Server{
		Addr:    ":" + listenAddr,
		Handler: router,
	}

	mapHandlers(router, service)

	return server
}
