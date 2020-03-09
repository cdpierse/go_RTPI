package models

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cdpierse/go_dublin_bus/constants"
	"github.com/gorilla/mux"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

// Appends stop specific endpoint to base server address
const (
	StopsURL = constants.RTPIBaseServer + "busstopinformation"
)

// Operator struct mapping fields for a
// service operator including the routes
// covered by that operator
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

// will cache this in future version and use better naming
func getAllStops() []Stop {
	body := GetRequestBody(StopsURL)
	stops := unpackStopResponseResults(body)
	return stops

}

func checkOperatorPresent(source, target string) bool {

	if strings.ToLower(source) == strings.ToLower(target) {
		return true
	}

	return false

}

// GetStopByQueryVals is effectively a combinaiton of all
// previous Stop GET requests that returns all results where
// a match is found for any given query key:value pair.
func filterByQuery(stops []Stop, r *http.Request) []Stop {
	var filteredStops []Stop
	stopParam := r.URL.Query().Get("stop_id")
	nameParam := r.URL.Query().Get("stop_name")
	operatorParam := r.URL.Query().Get("operator")

	if stopParam != "" {
		for _, item := range stops {
			if stopParam == item.Stopid {
				if operatorParam != "" {
					numOperators := len(item.Operators)
					for i := 0; i < numOperators; i++ {
						operatorName := item.Operators[i].Name
						if checkOperatorPresent(operatorParam, operatorName) {
							filteredStops = append(filteredStops, item)
						}
					}
				} else {
					filteredStops = append(filteredStops, item)
				}

			}
		}
	}

	if nameParam != "" {
		if len(filteredStops) != 0 {
			stops = filteredStops
		}
		for _, item := range stops {
			rankFullName := fuzzy.RankMatch(strings.ToLower(nameParam), strings.ToLower(item.Fullname))
			rankShortName := fuzzy.RankMatch(strings.ToLower(nameParam), strings.ToLower(item.Shortname))
			if (rankFullName > 7 && rankFullName <= 10) ||
				(rankShortName > 7 && rankShortName <= 10) {
				if operatorParam != "" {
					numOperators := len(item.Operators)
					for i := 0; i < numOperators; i++ {
						operatorName := item.Operators[i].Name
						if checkOperatorPresent(operatorParam, operatorName) {
							filteredStops = append(filteredStops, item)
						}
					}
				} else {
					filteredStops = append(filteredStops, item)
				}

			}
		}
	}

	return filteredStops
}

// GetStops returns all stops defined in the host system along
// with metadata for each stop returned.
func GetStops(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	queries := r.URL.Query()
	if len(queries) != 0 {
		stops = filterByQuery(stops, r)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stops)

}

// GetStop returns the metadata for a Stop with a given stop ID
// with metadata for each stop returned.
func GetStop(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	queries := r.URL.Query()
	if len(queries) != 0 {
		stops = filterByQuery(stops, r)
	}

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

	queries := r.URL.Query()
	if len(queries) != 0 {
		stops = filterByQuery(stops, r)
	}
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

// GetStopByFuzzyName will attempt to return all partial stop matches
// for a provided query string. It uses fuzzy searching to attempt to match
// a source string
func GetStopByFuzzyName(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	queries := r.URL.Query()
	if len(queries) != 0 {
		stops = filterByQuery(stops, r)
	}
	found := false

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range stops {
		rankFullName := fuzzy.RankMatch(strings.ToLower(params["stop_name"]), strings.ToLower(item.Fullname))
		rankShortName := fuzzy.RankMatch(strings.ToLower(params["stop_name"]), strings.ToLower(item.Shortname))
		if (rankFullName > 7 && rankFullName <= 20) ||
			(rankShortName > 7 && rankShortName <= 20) {
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

// GetStopsByOperator returns all stops that are serviced by
// the requested operator name i.e. BE (Bus Eireann).
func GetStopsByOperator(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	queries := r.URL.Query()
	if len(queries) != 0 {
		stops = filterByQuery(stops, r)
	}
	found := false
	filterByQuery(stops, r)

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, item := range stops {
		numOperators := len(item.Operators)
		for i := 0; i < numOperators; i++ {
			operatorName := item.Operators[i].Name
			if strings.ToLower(operatorName) == strings.ToLower(params["operator_name"]) {
				found = true
				json.NewEncoder(w).Encode(item)

			}
		}
	}
	if found {
		return
	} else {
		json.NewEncoder(w).Encode(&Stop{})
	}
}

// GetNearbyStopsByID returns all stops within distance n for a provided
// stop id
func GetNearbyStopsByID(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()

	found := false
	var sourcePoint orb.Point
	var targetPoint orb.Point
	maxDistance := 500.00
	queries := r.URL.Query()
	if len(queries) != 0 {
		md := r.URL.Query().Get("max_distance")
		if md != "" {
			maxDistance, _ = strconv.ParseFloat(md, 8)
		}

	}

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, item := range stops {
		if item.Stopid == params["stop_id"] {
			stopLat, err := strconv.ParseFloat(item.Latitude, 5)
			if err != nil {
				log.Panicln(err)
			}
			stopLong, err := strconv.ParseFloat(item.Longitude, 5)
			if err != nil {
				log.Panicln(err)
			}

			sourcePoint = orb.Point{stopLong, stopLat}
		}
	}
	for _, item := range stops {
		if item.Stopid != params["stop_id"] {
			targetLat, err := strconv.ParseFloat(item.Latitude, 5)
			if err != nil {
				log.Panicln(err)
			}
			targetLong, err := strconv.ParseFloat(item.Longitude, 5)
			if err != nil {
				log.Panicln(err)
			}
			targetPoint = orb.Point{targetLong, targetLat}
			if geo.Distance(sourcePoint, targetPoint) <= maxDistance {
				found = true
				json.NewEncoder(w).Encode(item)

			}
		}
	}
	if found {
		return
	} else {
		json.NewEncoder(w).Encode(&Stop{})
	}
}

func GetStopsByDistance(w http.ResponseWriter, r *http.Request) {
	stops := getAllStops()
	w.Header().Set("Content-Type", "application/json")

	found := false
	var sourcePoint orb.Point
	var targetPoint orb.Point
	maxDistance := 500.00

	queries := r.URL.Query()
	if len(queries) != 0 {

		lat, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 8)
		if err != nil {
			log.Panicln(err)
		}
		long, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 8)
		if err != nil {
			log.Panicln(err)
		}
		md := r.URL.Query().Get("max_distance")
		if md != "" {
			maxDistance, _ = strconv.ParseFloat(md, 8)
		}
		sourcePoint = orb.Point{long, lat}
	} else {
		log.Panicln("No query parameters for either lat or longitude found. ")
	}

	for _, item := range stops {
		targetLat, err := strconv.ParseFloat(item.Latitude, 5)
		if err != nil {
			log.Panicln(err)
		}
		targetLong, err := strconv.ParseFloat(item.Longitude, 5)
		if err != nil {
			log.Panicln(err)
		}
		targetPoint = orb.Point{targetLong, targetLat}

		if geo.Distance(sourcePoint, targetPoint) <= maxDistance {
			found = true
			json.NewEncoder(w).Encode(item)

		}

	}
	if found {
		return
	} else {
		json.NewEncoder(w).Encode(&Stop{})
	}


}
