package moov

import (
	"fmt"
	"time"
)

type Statement struct {
	StatementID                string              `json:"statementID,omitempty"`
	StatementName              string              `json:"statementName,omitempty"`
	FileName                   string              `json:"fileName,omitempty"`
	FileSize                   int64               `json:"fileSize,omitempty"`
	BillingPeriodStartDateTime time.Time           `json:"billingPeriodStartDateTime,omitempty"`
	BillingPeriodEndDateTime   time.Time           `json:"billingPeriodEndDateTime,omitempty"`
	SubscriptionIDs            []string            `json:"subscriptionIDs,omitempty"`
	CreatedOn                  time.Time           `json:"createdOn,omitempty"`
	UpdatedOn                  time.Time           `json:"updatedOn,omitempty"`
	Summary                    Summary             `json:"summary,omitempty"`
	CardAcquiringFees          CardAcquiringFees   `json:"cardAcquiringFees,omitempty,omitempty"`
	ACHFees                    *ACHFees            `json:"achFees,omitempty,omitempty"`
	InstantPaymentFees         *InstantPaymentFees `json:"instantPaymentFees,omitempty,omitempty"`
	PlatformFees               *PlatformFees       `json:"platformFees,omitempty,omitempty"`
	OtherCardFees              *OtherCardFees      `json:"otherCardFees,omitempty,omitempty"`
}

type Summary struct {
	CardAcquiring   SummaryDetailsWithInterchange `json:"cardAcquiring,omitempty"`
	ACH             *SummaryDetails               `json:"ach,omitempty"`
	InstantPayments *SummaryDetails               `json:"instantPayments,omitempty"`
	PlatformFees    *AmountDecimal                `json:"platformFees,omitempty"`
	AdjustmentFees  *AmountDecimal                `json:"adjustmentFees,omitempty"`
	OtherFees       *AmountDecimal                `json:"otherFees,omitempty"`
	Total           *AmountDecimal                `json:"total,omitempty"`
}

type SummaryDetails struct {
	FeeAmount *AmountDecimal `json:"feeAmount,omitempty"`
}

type SummaryDetailsWithInterchange struct {
	SummaryDetails
	InterchangeFees *InterchangeSummaryFees `json:"interchangeFees,omitempty"`
}

type InterchangeSummaryFees struct {
	Visa            AmountDecimal `json:"visa,omitempty"`
	Mastercard      AmountDecimal `json:"mastercard,omitempty"`
	Discover        AmountDecimal `json:"discover,omitempty"`
	AmericanExpress AmountDecimal `json:"americanExpress,omitempty"`
}

type CardAcquiringFees struct {
	AmericanExpress CardBrandFees `json:"americanExpress"`
	Discover        CardBrandFees `json:"discover"`
	Mastercard      CardBrandFees `json:"mastercard"`
	Visa            CardBrandFees `json:"visa"`
}

type CardBrandFees struct {
	CompletedDomesticInPerson      *CountAndAmount          `json:"completedDomesticInPerson,omitempty"`
	CompletedDomesticOnline        *CountAndAmount          `json:"completedDomesticOnline,omitempty"`
	CompletedInternationalInPerson *CountAndAmount          `json:"completedInternationalInPerson,omitempty"`
	CompletedInternationalOnline   *CountAndAmount          `json:"completedInternationalOnline,omitempty"`
	Declines                       *CountAndAmount          `json:"declines,omitempty"`
	Interchange                    *CountAndAmount          `json:"interchange,omitempty"`
	InterchangePrograms            *[]InterchangeProgramFee `json:"interchangePrograms,omitempty"`
	NetworkPassthrough             *CountAndAmount          `json:"networkPassthrough,omitempty"`
	Refunds                        *CountAndAmount          `json:"refunds,omitempty"`
	Total                          CountAndAmount           `json:"total,omitempty"`
}

type CountAndAmount struct {
	Amount AmountDecimal `json:"amount,omitempty"`
	Count  int64         `json:"count,omitempty"`
}

type InterchangeProgramFee struct {
	Count          int64         `json:"count,omitempty"`
	PerItemRate    AmountDecimal `json:"perItemRate,omitempty"`
	PercentageRate string        `json:"percentageRate,omitempty"`
	ProgramName    string        `json:"programName,omitempty"`
	TransferVolume AmountDecimal `json:"transferVolume,omitempty"`
	Total          AmountDecimal `json:"total,omitempty"`
}

type ACHFees struct {
	Debits             CountAndAmount `json:"debits,omitempty"`
	NoticeOfChange     CountAndAmount `json:"noticeOfChange,omitempty"`
	Return             CountAndAmount `json:"return,omitempty"`
	SameDayCredit      CountAndAmount `json:"sameDayCredit,omitempty"`
	StandardCredit     CountAndAmount `json:"standardCredit,omitempty"`
	UnauthorizedReturn CountAndAmount `json:"unauthorizedReturn,omitempty"`
	Total              CountAndAmount `json:"total,omitempty"`
}

type InstantPaymentFees struct {
	InstantVerification     CountAndAmount  `json:"instantVerification,omitempty"`
	PullFromCardDecline     *CountAndAmount `json:"pullFromCardDecline,omitempty"`
	PullFromCardRefund      *CountAndAmount `json:"pullFromCardRefund,omitempty"`
	PullFromCardTransaction CountAndAmount  `json:"pullFromCardTransaction,omitempty"`
	PushToCardDecline       *CountAndAmount `json:"pushToCardDecline,omitempty"`
	PushToCardTransaction   CountAndAmount  `json:"pushToCardTransaction,omitempty"`
	RtpCreditTransaction    CountAndAmount  `json:"rtpCreditTransaction,omitempty"`
	RtpDecline              *CountAndAmount `json:"rtpDecline,omitempty"`
	Total                   CountAndAmount  `json:"total,omitempty"`
}

type PlatformFees struct {
	MerchantPCIFee AmountDecimal `json:"merchantPCIFee,omitempty"`
	WalletFee      AmountDecimal `json:"walletFee,omitempty"`
	Total          AmountDecimal `json:"total,omitempty"`
}

type OtherCardFees struct {
	CardAccountUpdater CountAndAmount `json:"cardAccountUpdater,omitempty"`
	CardVerification   CountAndAmount `json:"cardVerification,omitempty"`
	Disputes           CountAndAmount `json:"disputes,omitempty"`
	NameVerification   CountAndAmount `json:"nameVerification,omitempty"`
	Total              CountAndAmount `json:"total,omitempty"`
}

type ListStatementFilter callArg

func WithStatementCount(c int) ListStatementFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithStatementSkip(c int) ListStatementFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithBillingPeriodStartDateTime(start time.Time) ListStatementFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["billingPeriodStartDateTime"] = start.Format(time.RFC3339)
		return nil
	})
}

func WithBillingPeriodEndDateTime(end time.Time) ListStatementFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["billingPeriodEndDateTime"] = end.Format(time.RFC3339)
		return nil
	})
}
