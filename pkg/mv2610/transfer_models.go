package mv2610

import (
	mapi "github.com/moovfinancial/moov-api-models/go/v20261000"
	"github.com/moovfinancial/moov-go/pkg/moov"
)

type AmountDecimal = mapi.AmountDecimal
type Capture = mapi.Capture
type CaptureStatus = mapi.CaptureStatus
type CardAcquiringRefund = mapi.CardAcquiringRefund
type CreateRefund = mapi.CreateRefund
type CreateRefundResponse = mapi.CreateRefundResponse
type CreateReversal = mapi.CreateReversal
type CreateTransfer = mapi.CreateTransfer
type CreateTransferACHAddendaRecord = mapi.CreateTransferACHAddendaRecord
type CreateTransferAmountDetails = mapi.CreateTransferAmountDetails
type CreateTransferDestination = mapi.CreateTransferDestination
type CreateTransferDestinationACH = mapi.CreateTransferDestinationACH
type CreateTransferDestinationCard = mapi.CreateTransferDestinationCard
type CreateTransferFacilitatorFee = mapi.CreateTransferFacilitatorFee
type CreateTransferLineItem = mapi.CreateTransferLineItem
type CreateTransferLineItemOption = mapi.CreateTransferLineItemOption
type CreateTransferLineItems = mapi.CreateTransferLineItems
type CreateTransferSource = mapi.CreateTransferSource
type CreateTransferSourceACH = mapi.CreateTransferSourceACH
type CreateTransferSourceCard = mapi.CreateTransferSourceCard
type CardPayoutType = mapi.CardPayoutType
type DebitHoldPeriod = mapi.DebitHoldPeriod
type RefundAmountDetails = mapi.RefundAmountDetails
type Reversal = mapi.Reversal
type ReversalAmountDetails = mapi.ReversalAmountDetails
type SECCode = mapi.SECCode
type Transfer = mapi.Transfer
type TransferAmountDetails = mapi.TransferAmountDetails
type TransferAuthorization = mapi.TransferAuthorization
type TransferDestination = mapi.TransferDestination
type TransferProcessingDetails = mapi.TransferProcessingDetails
type TransferRailOptions = mapi.TransferRailOptions
type TransferSource = mapi.TransferSource
type TransferStatus = mapi.TransferStatus
type TransferType = mapi.TransferType
type TransactionSource = mapi.TransactionSource

type TransferStarted = moov.TransferStarted
type CreateTransferOptions = moov.CreateTransferOptions
type CreateTransferOptionsTarget = moov.CreateTransferOptionsTarget
