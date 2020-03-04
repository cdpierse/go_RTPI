package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/cdpierse/go_dublin_bus/models"
)

func main() {
	fmt.Println("Starting Server...")
	router := mux.NewRouter()

	router.HandleFunc("/stops", models.GetStops).Methods("GET")
	router.HandleFunc("/stops/{stop_id}", models.GetStop).Methods("GET")


	http.ListenAndServe(":8000",router)
}
