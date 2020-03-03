package models

import (
	_"encoding/json"
	"fmt"

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

	fmt.Println(len(body))
	// json.Unmarshal(body,&stops)

	// for k,_ := range stops {
	// 	fmt.Printf("keys[%s] \n",k)
	// }
	
	return "hello"

}

func GetStop() {

}
