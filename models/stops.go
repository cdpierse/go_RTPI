package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/cdpierse/go_dublin_bus/constants"
	"github.com/gorilla/mux"
)

// Stop is a struct mapping out the fields returned
// from an RTPI response "results" field for Stop(s).
type Stop struct {
	Stopid             string `json:"stopid"`
	Displaystopid      string `json:"displaystopid"`
	Shortname          string `json:"shortname"`
	Shortnamelocalized string `json:"shortnamelocalized"`
	Fullname           string `json:"fullname"`
	Fullnamelocalized  string `json:"fullnamelocalized"`
	Latitude           string `json:"latitude"`
	Longitude          string `json:"longitude"`
	Lastupdated        string `json:"lastupdated"`
	Operators          []struct {
		Name         string   `json:"name"`
		Operatortype int      `json:"operatortype"`
		Routes       []string `json:"routes"`
	} `json:"operators"`
}

// StopsResponse struct mapping out RTPI JSON response object for stops.
// Huge thank you to Matt Holt (@mholt6) for developing
// https://mholt.github.io/json-to-go/ as that allowed me to generate this
// struct
type StopsResponse struct {
	Errorcode       string `json:"errorcode"`
	Errormessage    string `json:"errormessage"`
	Numberofresults int    `json:"numberofresults"`
	Timestamp       string `json:"timestamp"`
	Results         []Stop `json:"results"`
}

// Appends stop specific endpoint to base server address
const (
	StopsURL = constants.RTPIBaseServer + "busstopinformation"
)

func unpackStopResponseResults(responseJSON []byte) []Stop {
	var res StopsResponse

	err := json.Unmarshal(responseJSON, &res)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Succesfully retrieved All Stops")

	return res.Results

}

// GetStops returns all stops defined in the host system along 
// with metadata for each stop returned. 
func GetStops(w http.ResponseWriter, r *http.Request) {
	body := GetRequestBody(StopsURL)
	stops := unpackStopResponseResults(body)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stops)

}
// GetStops returns the metadata for a Stop with a given stop ID
// with metadata for each stop returned. 
func GetStop(w http.ResponseWriter, r *http.Request) {
	body := GetRequestBody(StopsURL)
	stops := unpackStopResponseResults(body)

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(("I get here at least"))
	params := mux.Vars(r)
	for _, item := range stops {
		if item.Stopid == params["stop_id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
			
	}
	json.NewEncoder(w).Encode(&Stop{})

	
	// nothing in here

}
