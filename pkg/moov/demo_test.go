package moov_test

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/moovfinancial/money/v2"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

/*
Problem:
Cancelling a transfer through the Moov Dashboard is great
if you just have to cancel one or a few transfers.

But what if you have A LOT of transfers to cancel?
It's not practical to cancel them one-by-one in the Dashboard.

Also, cancellation windows may be shorter due to settlement times.

This can be harder if transfers that should be canceled
are sprinkled between transfers that should NOT be canceled.


Solution:
Use our API -- or even better, use our SDK where we've alredy done all the plumbing for you
-- to cancel transfers in bulk.


To determine which transfers to cancel,
rather than validating manually in the Dashbaord -- you can use properties of the transfer
(or any other object in our API) to validate whether to cancel a certain transfer or not.

*/

const (
	partnerAccountID       = "db04bf9d-91f6-4206-ba38-6844636532ad"
	awesomeDronesAccountID = "e5a43c02-559a-49e4-a0f8-5b180ec30202"
	dbFilename             = "sqlite.txt"
	maxXfrAmt              = 20

	// total transfers created = 3  * numTransfersPerType
	numTransfersPerType = 3
)

/*
Simulates a scenario where lots of transfers are created, but only some need to be canceled.

Create 60 transfers:
- 20 card-to-wallet
- 20 bank-to-wallet
- 20 bank-to-bank

Each transfer has a random amount between $1 and $20.

🐛 I have a bug in here where:
there's a 33% chance the transfer amount will be multipled by 100.
e.g. so instead of a $14 transfer, it'll be $1400.
*/
func TestDemo_CreateTransfers(t *testing.T) {
	mc := NewTestClient(t)
	require := require.New(t)
	ctx := context.Background()

	// Generating transfer options for source=AwesomeDrones and destination=MoovFinancial (partner)
	options, err := mc.TransferOptions(ctx, moov.CreateTransferOptions{
		Source:      moov.CreateTransferOptionsTarget{AccountID: awesomeDronesAccountID},
		Destination: moov.CreateTransferOptionsTarget{AccountID: partnerAccountID},
		Amount:      moov.Amount{Currency: "USD", Value: int64(100 * maxXfrAmt)},
	})
	NoResponseError(t, err)

	var (
		cardPMID         string
		achDebitFundPMID string
		walletPMID       string
		achCreditStdPMID string
	)

	// Collecting payment method IDs from the transfer options
	for _, opt := range options.SourceOptions {
		if opt.PaymentMethodType == moov.PaymentMethodType_CardPayment {
			cardPMID = opt.PaymentMethodID
		} else if opt.PaymentMethodType == moov.PaymentMethodType_AchDebitFund {
			achDebitFundPMID = opt.PaymentMethodID
		}
	}

	for _, opt := range options.DestinationOptions {
		if opt.PaymentMethodType == moov.PaymentMethodType_MoovWallet {
			walletPMID = opt.PaymentMethodID
		} else if opt.PaymentMethodType == moov.PaymentMethodType_AchCreditStandard {
			achCreditStdPMID = opt.PaymentMethodID
		}
	}

	var (
		createdXfrIDs []string
		numBadXfrs    int
	)
	// A function to create a transfer
	create := func(srcPMID, dstPMID string, xfrType string) {
		randomNumber := rand.Intn(maxXfrAmt) + 1

		amt := moov.Amount{
			Currency: "USD",
			Value:    int64(100 * randomNumber),
		}
		hasBug := false
		randPercentage := rand.Float32()

		// 🐛🐛🐛🐛🐛🐛🐛🐛🐛🐛🐛
		if randPercentage <= .33333 {
			hasBug = true
			amt.Value *= 100
			numBadXfrs++
		}

		xfr, _, err := mc.CreateTransfer(ctx, partnerAccountID, moov.CreateTransfer{
			Source:      moov.CreateTransfer_Source{PaymentMethodID: srcPMID},
			Destination: moov.CreateTransfer_Destination{PaymentMethodID: dstPMID},
			Amount:      amt,
		}).WaitForRailResponse()
		NoResponseError(t, err)

		createdXfrIDs = append(createdXfrIDs, xfr.TransferID)

		numXfr := len(createdXfrIDs)
		if hasBug {
			fmt.Printf("(%d) 😨😨😨 Created %s %s transfer\n", numXfr, prettyPrintAmount(xfr.Amount), xfrType)
		} else {
			fmt.Printf("(%d) 💸 Created %s %s transfer\n", numXfr, prettyPrintAmount(xfr.Amount), xfrType)
		}

	}

	// Collect all transferIDs and writing the IDs to sqlite.txt file
	f, err := os.Create(dbFilename)
	require.NoError(err)
	defer f.Close()

	fmt.Printf("🤑 Creating transfers 🤑\n\n")

	for range numTransfersPerType {
		create(cardPMID, walletPMID, "card-to-wallet")
	}

	for range numTransfersPerType {
		create(achDebitFundPMID, walletPMID, "bank-to-wallet")
	}

	for range numTransfersPerType {
		create(achDebitFundPMID, achCreditStdPMID, "bank-to-bank")
	}

	b := strings.Join(createdXfrIDs, "\n")
	_, err = f.WriteString(b)
	require.NoError(err)

	// Print results
	fmt.Println()
	fmt.Println("Result: ")
	fmt.Printf("Created %d transfers successfully\n", len(createdXfrIDs))
	fmt.Printf("But %d out of %d transfers have WRONG amounts \n\n", numBadXfrs, len(createdXfrIDs))
}

/*
Cancels the transfers with invalid amounts:
1. Will fetch the transferIDs from sqlite.txt file
2. Determine if the transfer should be canceled (by checking if amount > $100)
3. Initiate the cancellation
4. Check the the cancellation completed
*/
func TestDemo_CancelTransfers(t *testing.T) {
	mc := NewTestClient(t)
	require := require.New(t)
	ctx := context.Background()

	f, err := os.Open(dbFilename)
	require.NoError(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	xfrIDsToAmount := map[string]moov.Amount{}
	for scanner.Scan() {
		xfrID := scanner.Text()
		if len(xfrID) == 0 {
			continue
		}

		xfr, err := mc.GetTransfer(ctx, partnerAccountID, xfrID)
		require.NoError(err)

		// Cancel any transfers greater than $100
		if xfr.Amount.Value > 100_00 {
			xfrIDsToAmount[xfr.TransferID] = xfr.Amount
		}
	}

	var numCancellations int

	fmt.Println("⚠️ Cancelling transfers ⚠️")
	fmt.Println()
	for xfrID, amount := range xfrIDsToAmount {
		cancellation, err := mc.CancelTransfer(ctx, partnerAccountID, xfrID)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			require.Equal(moov.CancellationStatus_Completed, cancellation.Status)
			fmt.Printf("❌ Canceled transfer with amount %s\n", prettyPrintAmount(amount))

			numCancellations += 1
		}
	}

	fmt.Printf("\n🎉 Canceled %d transfers\n\n\n", numCancellations)
}

func prettyPrintAmount(amt moov.Amount) string {
	moneyAmt, err := money.FromInt64(amt.Value, amt.Currency)
	if err != nil {
		panic(err)
	}

	return "$" + moneyAmt.ValueString()
}
