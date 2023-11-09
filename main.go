package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	moov "github.com/moovfinancial/moov-go/pkg"
)

func main() {

	// Setup Moov client
	// Generate API keys from the Moov Dashboard. Sandbox keys are not compatable with production.

	creds := moov.Credentials{
		AccountID: "638481a5-5205-406c-84c7-2fc2239105d1",
		PublicKey: "Qo0j0ChSalMKmRI_",
		SecretKey: "nbyCvG0r0MpM2gxkzo2u2G70jr6ZtFdm",
		Domain:    "localhost",
	}

	mc, err := moov.NewClient(creds)
	if err != nil {
		log.Fatal(err)
	}

	token, err := mc.BasicAuthToken()
	if err != nil {
		log.Fatal(err)
	}
	// Access token for website/mobile app
	log.Print(token.AccessToken)

	// client := moov.NewClient(creds)

	//accessToken, _ := client.BasicAuthToken()

	//log.Println(accessToken.AccessToken)

	// account := moov.Account{
	// 	AccountType: moov.INDIVIDUAL,
	// 	Profile: moov.Profile{
	// 		Individual: moov.Individual{
	// 			Name: moov.Name{
	// 				FirstName: "Wade",
	// 				LastName:  "Arnold",
	// 			},
	// 			Phone: moov.Phone{
	// 				Number:      "555-555-5555",
	// 				CountryCode: "1",
	// 			},
	// 			Address: moov.Address{
	// 				AddressLine1: "123 Main St",
	// 			},
	// 		},
	// 	},
	// }

	//accountJSON, _ := json.MarshalIndent(account, "", "  ")

	//log.Println(string(accountJSON))

	//valid response Body: {"mode":"sandbox","accountID":"638481a5-5205-406c-84c7-2fc2239105d1","accountType":"individual","displayName":"Wade Arnold","profile":{"individual":{"name":{"firstName":"Wade","lastName":"Arnold"},"phone":{"number":"5555555555","countryCode":"1"},"birthDateProvided":false,"governmentIDProvided":false}},"verification":{"verificationStatus":"unverified","status":"unverified"},"foreignID":"your-correlation-id","createdOn":"2023-11-08T23:06:16.168497001Z","updatedOn":"2023-11-08T23:06:16.168497001Z"}
	bytesCompare := []byte(`{"mode":"sandbox","accountID":"638481a5-5205-406c-84c7-2fc2239105d1","accountType":"individual","displayName":"Wade Arnold","profile":{"individual":{"name":{"firstName":"Wade","lastName":"Arnold"},"phone":{"number":"5555555555","countryCode":"1"},"birthDateProvided":false,"governmentIDProvided":false}},"verification":{"verificationStatus":"unverified","status":"unverified"},"foreignID":"your-correlation-id","createdOn":"2023-11-08T23:06:16.168497001Z","updatedOn":"2023-11-08T23:06:16.168497001Z"}`)
	//bytesCompare := []byte(`{"accountType":"individual","profile":{"individual":{"name":{"firstName":"Rick","lastName":"Arnold"},"phone":{"number":"5555555555","countryCode":"1"}}}}`)
	account := new(moov.Account)
	err = json.NewDecoder(bytes.NewReader(bytesCompare)).Decode(&account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", account)

	//bytejson := []byte(`{"accountType": "business",
	//			"profile": {"business": {"legalBusinessName": "Whole Body Fitness LLC","businessType": "llc",}},"foreignId": "your-correlation-id",
	//			"capabilities": ["transfers"]}`)

	// bytejsonIndividual := []byte(`{
	// 	"accountType": "individual",
	// 	"profile": {
	// 	  "individual": {
	// 		"name": {
	// 		  "firstName": "Wade",
	// 		  "lastName": "Arnold"
	// 		},
	// 		"phone": {
	// 		  "number": "555-555-5555",
	// 		  "countryCode": "1"
	// 		}
	// 	  }
	// 	},
	// 	"foreignId": "your-correlation-id"
	//   }`)

	/**

		req, err := http.NewRequest("POST", "https://api.moov.io/accounts", bytes.NewBuffer(bytesCompare))
		if err != nil {
			// Todo: return an error
			log.Fatal(err)
		}
		var PTransport http.RoundTripper = &http.Transport{Proxy: http.ProxyFromEnvironment}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(client.KeyPublic, client.KeySecret)
		httpC := &http.Client{Transport: PTransport}
		resp, err := httpC.Do(req)
		if err != nil {
			// Todo: return an error
			log.Fatal(err)
		}

		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

	**/

	/**
		req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.KeyPublic, c.KeySecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	if err := json.Unmarshal(resBody, &token); err != nil { // Parse []byte to go struct pointer
		// Todo: return an error
		log.Fatal("Can not unmarshal JSON")
	}

	**/

	/*
	   *

	   	account := new(moov.Account)
	   	account.accountType = INDIVIDUAL

	   	individual := Individual{}
	*/
}
