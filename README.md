# moov-go
Golang client for interacting with Moov API server side 

### Installation 

```bash 
# Go Client 
go get github.com/moovfinancial/moov-go/pkg
```

## Basic Usage 

```go 
package main

import (
	"log"

	moov "github.com/moovfinancial/moov-go/pkg"
)

func main() {

	// Setup Moov client
	// Generate API keys from the Moov Dashboard. Sandbox keys are not compatable with production.

	creds := moov.Credentials{
		AccountID: "638481a5-FAKE-406c-84c7-2fc2239105d1",
		PublicKey: "Qo0j0ChFAKEKmRI_",
		SecretKey: "nbyFAKEr0MpM2gxkzo2u2G70jr6ZtFdm",
		Domain:    "localhost",
	}

	mc := moov.NewClient(creds)

	token, err := mc.BasicAuthToken()
	if err != nil {
		log.Fatal(err)
	}
	// Access token for website/mobile app to instantiate Moov.js
	log.Print(token.AccessToken)

```
