package moov_test

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

const (
	partnerAccountID       = "db04bf9d-91f6-4206-ba38-6844636532ad"
	awesomeDronesAccountID = "e5a43c02-559a-49e4-a0f8-5b180ec30202"
	dbFilename             = "sqlite.txt"
)

func TestX(t *testing.T) {
	mc := NewTestClient(t)
	require := require.New(t)
	ctx := context.Background()

	max := 20

	options, err := mc.TransferOptions(BgCtx(), moov.CreateTransferOptions{
		Source:      moov.CreateTransferOptionsTarget{AccountID: awesomeDronesAccountID},
		Destination: moov.CreateTransferOptionsTarget{AccountID: partnerAccountID},
		Amount:      moov.Amount{Currency: "USD", Value: int64(100 * max)},
	})
	NoResponseError(t, err)

	var (
		// pmIds
		cardPMID         string
		achDebitFundPMID string

		//
		walletPMID       string
		achCreditStdPMID string
	)
	fmt.Println(achDebitFundPMID, achCreditStdPMID)

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

	var createdXfrIDs []string
	create := func(srcPMID, dstPMID string) {
		randomNumber := rand.Intn(max) + 1
		randPercentage := rand.Float32() * .33

		amt := moov.Amount{
			Currency: "USD",
			Value:    int64(100 * randomNumber),
		}
		if randPercentage <= .33 {
			amt.Value *= 100
		}

		completed, _, err := mc.CreateTransfer(BgCtx(), partnerAccountID, moov.CreateTransfer{
			Source:      moov.CreateTransfer_Source{PaymentMethodID: srcPMID},
			Destination: moov.CreateTransfer_Destination{PaymentMethodID: dstPMID},
			Amount:      amt,
		}).WaitForRailResponse()
		NoResponseError(t, err)

		createdXfrIDs = append(createdXfrIDs, completed.TransferID)

	}

	f, err := os.Create(dbFilename)
	require.NoError(err)
	defer f.Close()

	b := strings.Join(createdXfrIDs, "\n")
	_, err = f.WriteString(b)
	require.NoError(err)

	numTransfersPerType := 1

	for range numTransfersPerType {
		create(cardPMID, walletPMID) // card-to-wallet
	}

	for range numTransfersPerType {
		create(achDebitFundPMID, walletPMID) // bank-to-wallet
	}

	for range numTransfersPerType {
		create(achDebitFundPMID, achCreditStdPMID) // bank-to-bank
	}

	f, err = os.Open(dbFilename)
	require.NoError(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var xfrIDs []string
	for scanner.Scan() {
		xfrID := scanner.Text()
		if len(xfrID) == 0 {
			continue
		}

		xfr, err := mc.GetTransfer(ctx, partnerapartnerAccountID, xfrID)
		require.NoError(err, )


		createdXfrIDs = append(createdXfrIDs, xfrID)
	}

	for _, xfrID := range createdXfrIDs {
		cancellation, err := mc.CancelTransfer(BgCtx(), partnerAccountID, xfrID)
		NoResponseError(t, err)

		require.Equal(moov.CancellationStatus_Completed, cancellation.Status)
	}

	// t.Run("create and get cancellation", func(t *testing.T) {
	// 	createdCancellation, err := mc.CancelTransfer(BgCtx(), FACILITATOR_ID, transferID)
	// 	NoResponseError(t, err)
	// 	require.NotEmpty(t, createdCancellation)

	// 	fetchedCancellation, err := mc.GetCancellation(BgCtx(), FACILITATOR_ID, transferID, createdCancellation.CancellationID)
	// 	NoResponseError(t, err)
	// 	require.Equal(t, createdCancellation.CancellationID, fetchedCancellation.CancellationID)
	// })
}
