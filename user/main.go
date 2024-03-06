package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"server/db"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func identifierUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post ...")
	var p db.UserLogin
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	identifier := db.SelectUser(p)
	fmt.Println(identifier)
	if !identifier {
		http.Error(w, "Change name or/and password", http.StatusUnauthorized)
		return
	}
	io.WriteString(w, "Connected\n")
	w.WriteHeader(200)
	return
}

func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post ...")
	var p db.User
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db.InsertUser(p)
}

func isOK(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func main() {
	db.Database()
	myRouter := mux.NewRouter().StrictSlash(true)
	internRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/register", addUser).Methods(http.MethodPost)
	myRouter.HandleFunc("/login", identifierUser).Methods(http.MethodPost)
	internRouter.HandleFunc("/ping", isOK).Methods(http.MethodGet)
	internRouter.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	go func() {
		log.Fatal(http.ListenAndServe(":8001", internRouter))
	}()
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
