package transfers

import (
	"context"
	"fmt"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func ExampleCancelTransfer() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Errorf("new Moov client: %w", err)
		return
	}

	var (
		ctx              = context.Background()
		partnerAccountID = "00000000-00000000-00000000-00000000"
		transferID       = "00000000-00000000-00000000-00000000"
	)

	cancellation, err := mc.CancelTransfer(ctx, partnerAccountID, transferID)
	if err != nil {
		fmt.Printf("cancelling transfer: %w", err)
		return
	}

	fmt.Printf("Created cancellation: %+v", cancellation)
}

func ExampleGetCancellation() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Errorf("new Moov client: %w", err)
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
		fmt.Printf("getting cancellation: %w", err)
		return
	}

	fmt.Printf("Got cancellation: %+v", cancellation)
}
