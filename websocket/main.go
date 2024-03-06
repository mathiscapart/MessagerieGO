package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
)

func isOK(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	hub := newHub()
	go hub.run()
	myRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	internRouter := mux.NewRouter().StrictSlash(true)
	internRouter.HandleFunc("/ping", isOK).Methods(http.MethodGet)
	internRouter.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	go func() {
		log.Fatal(http.ListenAndServe(":8003", internRouter))
	}()
	log.Fatal(http.ListenAndServe(":8002", myRouter))
}
