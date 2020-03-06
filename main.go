package main

import (
	_ "log"
	_ "net/http"

	"github.com/cdpierse/go_dublin_bus/api"
)

func main() {
	api.InitializeAPI()
}
