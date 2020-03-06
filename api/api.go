package api

import (
	"github.com/cdpierse/go_dublin_bus/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// InitializeAPI instantiates a new router instance of
// our RTPI service, opening up its endpoints.
func InitializeAPI() {
	log.Println("Starting Up Server...")
	router := mux.NewRouter()
	log.Println("Server Online")

	router.HandleFunc("/stops", models.GetStops).Methods("GET")
	router.HandleFunc("/stops/id/{stop_id}", models.GetStop).Methods("GET")
	router.HandleFunc("/stops/stop_name/{stop_name}", models.GetStopByName).Methods("GET")
	router.HandleFunc("/stops/stop_name_fuzzy/{stop_name}", models.GetStopByFuzzyName).Methods("GET")
	router.HandleFunc("/stops/operator/{operator_name}", models.GetStopsByOperator).Methods("GET")

	http.ListenAndServe(":8000", router)

}
