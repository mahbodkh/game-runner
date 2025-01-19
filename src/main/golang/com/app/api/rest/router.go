package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404 Not Found")
	fmt.Println("Endpoint Hit: notFound")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("*", notFound)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}
