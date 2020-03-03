package models

import (
	"encoding/json"
	_ "fmt"

	"github.com/cdpierse/go_dublin_bus/constants"
	_ "github.com/gorilla/mux"
	// "io/ioutil"
	// "log"
	// "net/http"
)

type Stop struct {
	StopID   int    `json:"stop_id"`
	StopName string `json:"stop_name"`
	StopNameLocalized string `json:"stop_name_localized"`
	Operator string `json:"operator"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`

}

// Appends stop specific endpoint to base server address
const (
	StopsURL = constants.RTPIBaseServer + "busstopinformation"
)

func GetStops() string {

	body := GetRequestBody(StopsURL)
	
	return body

}
