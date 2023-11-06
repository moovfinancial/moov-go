package main

import (
	"log"

	moov "github.com/moovfinancial/moov-go/client"
)

func main() {

	// Setup Moov client
	// Generate API keys from the Moov Dashboard. Sandbox keys are not compatable with production.
	moov := new(moov.Client)
	moov.KeyPublic = "sandboxPublic"
	moov.KeySecret = "sandboxSecret"

	accessToken, _ := moov.BasicAuthToken()

	log.Println(accessToken.AccessToken)

}
