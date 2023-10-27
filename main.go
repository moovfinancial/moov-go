package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ClientCredentialsGrantToAccessTokenResponse struct {
	AccessToken  string `json:"access_token`
	RefreshToken string `json:"access_token`
	TokenType    string `json:"token_type`
	ExpiresIn    int32  `json:"expires_in`
	Scope        string `json:"scope"`
}

func main() {

	// Setup Moov client
	// Generate API keys from the Moov Dashboard. Sandbox keys are not compatable with production.

	KeyPublic := "Qo0j0ChSalMKmRI_"
	KeySecret := "nbyCvG0r0MpM2gxkzo2u2G70jr6ZtFdm"

	key := b64.StdEncoding.EncodeToString([]byte(KeyPublic + ":" + KeySecret))
	log.Println(key)

	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("scope", "/accounts.write")

	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(KeyPublic, KeySecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", resBody)

	tokenResponse := ClientCredentialsGrantToAccessTokenResponse{}
	if err := json.Unmarshal(resBody, &tokenResponse); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	log.Printf("tokenResponse: %s", tokenResponse.AccessToken)

}
