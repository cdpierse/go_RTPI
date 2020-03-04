package models

import (
	_ "bytes"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/cdpierse/go_dublin_bus/constants"
	_ "github.com/gorilla/mux"
)

// StopResult is a struct mapping out the fields returned
// in the results for Stop(s) GET request in the RTPI server.
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
	Errorcode       string       `json:"errorcode"`
	Errormessage    string       `json:"errormessage"`
	Numberofresults int          `json:"numberofresults"`
	Timestamp       string       `json:"timestamp"`
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
	fmt.Println(reflect.TypeOf(res))
	fmt.Println(reflect.TypeOf(res.Results))

	return res.Results

}

func GetStops() string {

	body := GetRequestBody(StopsURL)
	results := unpackStopResponseResults(body)
	fmt.Println(results[0:10])

	return "hello"

}

func GetStop() {
	// nothing in here

}
