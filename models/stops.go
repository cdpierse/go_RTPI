package models


import (
	"github.com/cdpierse/go_dublin_bus/constants"
	_"fmt"
)



type Stop struct {
	StopID int `json:"stop_id"`
	StopName string `json:"stop_name"`
	Operator string `json:"operator"`

}

// appends stop specific endpoint to base server address
const (
	StopsURL = constants.RTPIBaseServer + "busstopinformation"
)


func PrintBaseServer() string {
	return StopsURL
}