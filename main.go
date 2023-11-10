package main

import (
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

	token, err := mc.SingleUseAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	// Access token for website/mobile app
	log.Print(token.AccessToken)

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

	/*
	   	params := url.Values{}
	   	params.Add("grant_type", "client_credentials")
	   	params.Add("scope", "/accounts.write")

	   	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)

	   	func (c Client) Ping() {
	   	log.Println("ping")
	   	req, _ := http.NewRequest("POST", "https://api.moov.io/ping", nil)
	   	req.Header.Set("Accept", "application/json")
	   	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	   	//req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	   	client := &http.Client{}
	   	resp, err := client.Do(req)
	   	if err != nil {
	   		// Todo: return an error
	   		log.Fatal(err)
	   	}
	   	defer resp.Body.Close()

	   	body, err := io.ReadAll(resp.Body)
	   	if err != nil {
	   		// Todo: return an error
	   		log.Fatal(err)
	   	}

	   	fmt.Println("response Status:", resp.Status)
	   	fmt.Println("response Headers:", resp.Header)
	   	fmt.Println("response Body:", string(body))

	   }
	*/

}
