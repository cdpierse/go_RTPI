package models

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/cdpierse/go_dublin_bus/constants"
	"github.com/gorilla/mux"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Appends stop specific endpoint to base server address
const (
	StopsURL = constants.RTPIBaseServer + "busstopinformation"
)

type Operator struct {
	Name         string   `json:"name"`
	Operatortype int      `json:"operatortype"`
	Routes       []string `json:"routes"`
}

// Stop is a struct mapping out the fields returned
// from an RTPI response "results" field for Stop(s).
type Stop struct {
	Stopid             string     `json:"stopid"`
	Displaystopid      string     `json:"displaystopid"`
	Shortname          string     `json:"shortname"`
	Shortnamelocalized string     `json:"shortnamelocalized"`
	Fullname           string     `json:"fullname"`
	Fullnamelocalized  string     `json:"fullnamelocalized"`
	Latitude           string     `json:"latitude"`
	Longitude          string     `json:"longitude"`
	Lastupdated        string     `json:"lastupdated"`
	Operators          []Operator `json:"operators"`
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

func unpackStopResponseResults(responseJSON []byte) []Stop {
	var res StopsResponse

	err := json.Unmarshal(responseJSON, &res)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Succesfully retrieved All Stops")

	return res.Results

}
func getAllStops() []Stop {
	body := GetRequestBody(StopsURL)
	stops := unpackStopResponseResults(body)
	return stops

}

// GetStops returns all stops defined in the host system along
// with metadata for each stop returned.
func GetStops(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stops)

}

// GetStop returns the metadata for a Stop with a given stop ID
// with metadata for each stop returned.
func GetStop(w http.ResponseWriter, r *http.Request) {
	body := GetRequestBody(StopsURL)
	stops := unpackStopResponseResults(body)
	found := false

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range stops {
		if item.Stopid == params["stop_id"] {
			found = true
			json.NewEncoder(w).Encode(item)
			return
		}

	}
	if found {
		return
	} else {
		// if nothing is returned we can at least
		// write the structure of a stop to stream.
		json.NewEncoder(w).Encode(&Stop{})
	}

}

// GetStopByName returns all stop(s) where the query parameter
// stop_name is either equal to the the full name or short name
// of a stop.
func GetStopByName(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	found := false

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, item := range stops {
		if strings.ToLower(item.Fullname) == strings.ToLower(params["stop_name"]) ||
			strings.ToLower(item.Shortname) == strings.ToLower(params["stop_name"]) {
			found = true
			json.NewEncoder(w).Encode(item)

		}

	}
	if found {
		return
	} else {
		// if nothing is returned we can at least
		// write the structure of a stop to stream.
		json.NewEncoder(w).Encode(&Stop{})
	}

}

func GetStopByOperator(w http.ResponseWriter, r *http.Request) {
	// stops := getAllStops()
	// found := false

	// w.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(r)

}
// GetStopByFuzzyName will attempt to return all partial stop matches
// for a provided query string. It uses fuzzy searching to attempt to match
// a source string
func GetStopByFuzzyName(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	found := false
	log.Println(r.URL.Query())
	log.Println(mux.Vars(r))
	query := r.URL.Query()
	log.Println(query["key2"][0])

	

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range stops {
		rankFullName := fuzzy.RankMatch(strings.ToLower(params["stop_name"]), strings.ToLower(item.Fullname))
		rankShortName := fuzzy.RankMatch(strings.ToLower(params["stop_name"]), strings.ToLower(item.Shortname))
		if  (rankFullName > 7 && rankFullName <= 10) ||
			(rankShortName > 7 && rankShortName <= 10) {
			found = true
			json.NewEncoder(w).Encode(item)

		}

	}
	if found {
		return
	} else {
		// if nothing is returned we can at least
		// write the structure of a stop to stream.
		json.NewEncoder(w).Encode(&Stop{})
	}

}

// GetStopByQueryVals is effectively a combinaiton of all 
// previous Stop GET requests that returns all results where
// a match is found for any given query key:value pair. 
func GetStopByQueryVals(w http.ResponseWriter, r *http.Request) {
	// stops := getAllStops()
	// found := false
	// log.Println(r.URL.Query())

}
