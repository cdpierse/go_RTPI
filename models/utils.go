package models


import (
	"io/ioutil"
	"log"
	"net/http"
)

// GetRequestBody is a generic function for 
// wrapping GET requests in and return the 
// body of that request if succesfull
func GetRequestBody(url string) string {

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

	return string(body)

}