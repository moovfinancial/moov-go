package transfers

import (
	"context"
	"fmt"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func ExampleClient_CancelTransfer() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Errorf("new Moov client: %v", err)
		return
	}

	var (
		ctx              = context.Background()
		partnerAccountID = "00000000-00000000-00000000-00000000"
		transferID       = "00000000-00000000-00000000-00000000"
	)

	cancellation, err := mc.CancelTransfer(ctx, partnerAccountID, transferID)
	if err != nil {
		fmt.Printf("cancelling transfer: %v", err)
		return
	}

	fmt.Printf("Created cancellation: %+v", cancellation)
}

func ExampleClient_GetCancellation() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Errorf("new Moov client: %v", err)
		return
	}

	var (
		ctx              = context.Background()
		partnerAccountID = "00000000-00000000-00000000-00000000"
		transferID       = "00000000-00000000-00000000-00000000"
		cancellationID   = "00000000-00000000-00000000-00000000"
	)

	cancellation, err := mc.GetCancellation(ctx, partnerAccountID, transferID, cancellationID)
	if err != nil {
		fmt.Printf("getting cancellation: %v", err)
		return
	}

	fmt.Printf("Got cancellation: %+v", cancellation)
}
