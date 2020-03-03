package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_"reflect"
)

type RTPIResponse struct {
	Results []interface{} `json:"results"`
}

// Huge thank you to Matt Holt (@mholt6) for developing
// https://mholt.github.io/json-to-go/ as that allowed me to generate this 
// struct
type StopsResponse struct {
	Errorcode       string `json:"errorcode"`
	Errormessage    string `json:"errormessage"`
	Numberofresults int    `json:"numberofresults"`
	Timestamp       string `json:"timestamp"`
	Results         []struct {
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
	} `json:"results"`
}

// GetRequestBody is a generic function for
// wrapping GET requests in and return the
// body of that request if succesfull
func GetRequestBody(url string) []byte {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// makes sure that this closes after
	// function executes
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)

	}

	unpackResponseResults(body)

	return body

}

func unpackResponseResults(responseJSON []byte) {
	var sr StopsResponse
	// var results interface{}

	err := json.Unmarshal(responseJSON, &sr)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sr.Results[1].Stopid)
	// m := f.(map[string]interface{})

	// results := m["results"]
	// fmt.Println(len(results))

	// for k,v := range m {
	// 	fmt.Println(k)
	// 	if k == "numberofresults"{
	// 		fmt.Println(v)
	// 	}

	// 	if k == "results"{
	// 		fmt.Println(v)
	// 	}
	// }


}
