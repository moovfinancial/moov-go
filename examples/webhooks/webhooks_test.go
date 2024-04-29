package webhooks

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moovfinancial/moov-go/pkg/mhooks"
)

// Example handler func for processing a single event type (AccountCreated)
func ExampleHandler_SingleEvent() {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		secret := "your-webhook-signing-secret" // fetched from secure storage
		event, err := mhooks.ParseEvent(r, secret)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		accountCreated, err := event.AccountCreated()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		fmt.Printf("Got AccountCreated webhook with accountID=%v\n", accountCreated.AccountID)
		w.WriteHeader(200)
	}

	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}

// Example handler func for processing multiple event types (TransferCreated or TransferUpdated)
func ExampleHandler_MultipleEvents() {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		secret := "your-webhook-signing-secret" // fetched from secure storage
		event, err := mhooks.ParseEvent(r, secret)
		if err != nil {
			w.WriteHeader(500)
		}

		//nolint:exhaustive
		switch event.EventType {
		case mhooks.EventTypeTransferCreated:
			got, err := event.TransferCreated()
			if err != nil {
				w.WriteHeader(500)
				return
			}
			fmt.Printf("Got TransferCreated webhook with transferID=%v\n", got.TransferID)
		case mhooks.EventTypeTransferUpdated:
			got, err := event.TransferUpdated()
			if err != nil {
				w.WriteHeader(500)
				return
			}
			fmt.Printf("Got TransferUpdated webhook with transferID=%v\n", got.TransferID)
		}

		w.WriteHeader(200)
	}

	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}
