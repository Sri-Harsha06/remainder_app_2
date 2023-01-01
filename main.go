package main

import (
	"fmt"
	"net/http"
	"remainder_app_2/controller"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting the application...")
	router := mux.NewRouter()
	router.HandleFunc("/tmrevent", controller.Findtmrevents).Methods("GET")
	http.ListenAndServe(":12345", router)
}
