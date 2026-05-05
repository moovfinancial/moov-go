package transfers

import (
	"context"
	"fmt"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func ExampleClient_CreateTransferConfig() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Printf("new Moov client: %v\n", err)
		return
	}

	config, err := mc.CreateTransferConfig(context.Background(),
		"00000000-00000000-00000000-00000000",
		moov.UpsertTransferConfig{
			TipPresets: &moov.UpsertTipPresets{
				CalculationBasis:  moov.PtrOf(moov.TipCalculationBasis_PreTax),
				PercentageOptions: []int{10, 15, 20},
			},
		},
	)
	if err != nil {
		fmt.Printf("creating transfer config: %v\n", err)
		return
	}

	fmt.Printf("Created transfer config: %+v", config)
}

func ExampleClient_GetTransferConfig() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Printf("new Moov client: %v\n", err)
		return
	}

	config, err := mc.GetTransferConfig(context.Background(), "00000000-00000000-00000000-00000000")
	if err != nil {
		fmt.Printf("getting transfer config: %v\n", err)
		return
	}

	fmt.Printf("Got transfer config: %+v", config)
}

func ExampleClient_UpdateTransferConfig() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Printf("new Moov client: %v\n", err)
		return
	}

	config, err := mc.UpdateTransferConfig(context.Background(),
		"00000000-00000000-00000000-00000000",
		moov.UpsertTransferConfig{
			TipPresets: &moov.UpsertTipPresets{
				FixedAmountOptions: []moov.AmountDecimal{
					{Currency: "USD", ValueDecimal: "1.00"},
					{Currency: "USD", ValueDecimal: "2.00"},
					{Currency: "USD", ValueDecimal: "5.00"},
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("updating transfer config: %v\n", err)
		return
	}

	fmt.Printf("Updated transfer config: %+v", config)
}
