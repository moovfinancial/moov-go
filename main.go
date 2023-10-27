package main

import (
	"log"
	moov "moov-golang/client"
)

func main() {

	// Setup Moov client
	// Generate API keys from the Moov Dashboard. Sandbox keys are not compatable with production.
	moov := new(moov.Client)
	moov.KeyPublic = "Qo0j0ChSalMKmRI_"
	moov.KeySecret = "nbyCvG0r0MpM2gxkzo2u2G70jr6ZtFdm"

	accessToken, _ := moov.BasicAuthToken()

	log.Println(accessToken.AccessToken)

}
