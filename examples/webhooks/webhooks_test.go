package webhooks

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/moovfinancial/moov-go/pkg/mhooks"
)

// Example handler func for processing a single event type (AccountCreated)
//
// The handler will only set a non-200 status if an unexpected error occurs, and
// we want the webhook to be retried.
func ExampleHandler_SingleEvent() {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		secret := "your-webhook-signing-secret" // fetched from secure storage

		event, err := mhooks.ParseEvent(r, secret)
		if err != nil {
			msgPrefix := "parsing event: "
			if errors.Is(err, mhooks.ErrInvalidSignature) {
				msgPrefix += "invalid signature: "
			}
			fmt.Printf("%s %v\n", msgPrefix, err)
			w.WriteHeader(500)
			return
		}

		accountCreated, err := event.AccountCreated()
		if err != nil {
			fmt.Printf("getting AccountCreated from event: %v\n", err)
			w.WriteHeader(200)
			return
		}

		fmt.Printf("Got AccountCreated webhook with accountID=%v\n", accountCreated.AccountID)
		w.WriteHeader(200)
	}

	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}

// Example handler func for processing multiple event types (TransferCreated or TransferUpdated)
//
// The handler will only set a non-200 status if an unexpected error occurs, and
// we want the webhook to be retried.
func ExampleHandler_MultipleEvents() {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		secret := "your-webhook-signing-secret" // fetched from secure storage
		event, err := mhooks.ParseEvent(r, secret)
		if err != nil {
			msgPrefix := "parsing event: "
			if errors.Is(err, mhooks.ErrInvalidSignature) {
				msgPrefix += "invalid signature: "
			}
			fmt.Printf("%s %v\n", msgPrefix, err)
			w.WriteHeader(500)
			return
		}

		//nolint:exhaustive
		switch event.EventType {
		case mhooks.EventTypeTransferCreated:
			got, err := event.TransferCreated()
			if err != nil {
				fmt.Printf("getting TransferCreated from event: %v\n", err)
				w.WriteHeader(200)
				return
			}
			fmt.Printf("Got TransferCreated webhook with transferID=%v\n", got.TransferID)
		case mhooks.EventTypeTransferUpdated:
			got, err := event.TransferUpdated()
			if err != nil {
				fmt.Printf("getting TransferUpdated from event: %v\n", err)
				w.WriteHeader(200)
				return
			}
			fmt.Printf("Got TransferUpdated webhook with transferID=%v\n", got.TransferID)
		}

		w.WriteHeader(200)
	}

	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}
