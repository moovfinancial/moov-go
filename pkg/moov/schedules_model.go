package moov

import "time"

type Schedule struct {
	// prod or sandbox
	Mode string `json:"mode,omitempty"`

	ScheduleID string `json:"scheduleID,omitempty"`

	// This is the account ID of the source transfers that the transfer will run using.
	SourceAccountID string `json:"sourceAccountID,omitempty"`

	// This is the destination ID of the destination transfers that the transfer will run using.
	DestinationAccountID string `json:"destinationAccountID,omitempty"`

	// This is the partner ID that the transfer will run using.
	PartnerAccountID string `json:"partnerAccountID,omitempty"`

	// AccountID of the account that created it and is allowed to update it.
	OwnerID string `json:"ownerID,omitempty"`

	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	RecurTransfer *RecurTransfer `json:"recurTransfer,omitempty"`

	// List of all generated and manually added transfers to be made.
	Occurrences []TransferOccurrence `json:"occurrences,omitempty"`

	CreatedOn  time.Time  `json:"createdOn,omitempty"`
	UpdatedOn  time.Time  `json:"updatedOn,omitempty"`
	DisabledOn *time.Time `json:"disabledOn,omitempty" `
}

func (s Schedule) ToUpdateSchedule() UpdateSchedule {
	upsOccs := make([]UpdateTransferOccurrence, len(s.Occurrences))
	for i, occ := range s.Occurrences {
		upsOccs[i] = UpdateTransferOccurrence{
			OccurrenceID: &occ.OccurrenceID,
			Transfer:     occ.Transfer,
			RunOn:        occ.RunOn,
			Cancelled:    nil,
		}
	}

	return UpdateSchedule{
		Description:   s.Description,
		RecurTransfer: s.RecurTransfer,
		Occurrences:   upsOccs,
	}
}

// https://www.rfc-editor.org/rfc/rfc5545#section-3.3.10
type RecurTransfer struct {
	// Transfer values to use to create the transfer based on the recurRule
	// When changed, should just modify the transfer of the schedules
	Transfer ScheduleTransfer `json:"transfer,omitempty"`

	// This is the recurrence rule that is used to generate occurrences.
	// Generator available here: https://jkbrzt.github.io/rrule/
	// You can read the details of the format here: https://www.rfc-editor.org/rfc/rfc5545#section-3.3.10
	RecurrenceRule string `json:"recurrenceRule,omitempty"`

	// If the recurrence rule ends up being indefinite
	Indefinite bool `json:"indefinite,omitempty"`
}

type TransferOccurrence struct {
	ScheduleID string `json:"scheduleID,omitempty"`

	// Unique ID for updating a specific occurrence
	OccurrenceID string `json:"occurrenceID,omitempty"`

	// Mode to run the occurence under
	Mode string `json:"mode,omitempty"`

	// Transfer details that will be used.
	Transfer ScheduleTransfer `json:"transfer,omitempty"`

	// If this scheduled transfer was generated or manually added for say a correction
	// If a new interval is specified, all un-ran generated transfers will be re-generated
	Generated bool `json:"generated,omitempty"`

	// Flag if this is part of an indefinite schedule
	Indefinite bool `json:"indefinite,omitempty"`

	// Modified since generated. This could be switching just a single payment method
	Modified bool `json:"modified,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`

	// When the transfer was kicked off. If nil, hasn't ran. Normalize to UTC.
	RanOn *time.Time `json:"ranOn,omitempty"`

	// Ability to cancel this specific transfer from running
	CancelledOn *time.Time `json:"cancelledOn,omitempty"`

	// ID of the transfer that ran
	TransferID     *string `json:"transferID,omitempty"`
	TransferStatus *string `json:"transferStatus,omitempty"`
}

type CreateSchedule struct {
	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	RecurTransfer *RecurTransfer `json:"recurTransfer,omitempty"`

	// On creating the schedule we can use these occurrences as they planned the schedule
	Occurrences []CreateTransferOccurrence `json:"occurrences,omitempty"`
}

type CreateTransferOccurrence struct {
	// Transfer details that will be used.
	Transfer ScheduleTransfer `json:"transfer,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`
}

type UpdateSchedule struct {
	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	RecurTransfer *RecurTransfer `json:"recurTransfer,omitempty"`

	// On creating the schedule we can use these occurrences as they planned the schedule
	Occurrences []UpdateTransferOccurrence `json:"occurrences,omitempty"`
}

type UpdateTransferOccurrence struct {
	// Leave empty to add a new occurrence or set to the ID of the occurrence to change.
	OccurrenceID *string `json:"occurrenceID,omitempty"`

	// Transfer details that will be used.
	Transfer ScheduleTransfer `json:"transfer,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`

	// If nil, cancelledOn will be unchanged. If set true, it will be cancelled. If set false and hasn't ran yet will be uncancelled
	Cancelled *bool `json:"cancelled,omitempty"`
}

type ScheduleTransfer struct {
	Description string         `json:"description,omitempty"`
	Amount      ScheduleAmount `json:"amount,omitempty"`

	PartnerID   string                `json:"partnerAccountID,omitempty"`
	Source      SchedulePaymentMethod `json:"source,omitempty"`
	Destination SchedulePaymentMethod `json:"destination,omitempty"`
}

type ScheduleAmount struct {
	Value    int64  `json:"value,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type SchedulePaymentMethod struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty"`

	AchDetails  *ScheduleAchDetails  `json:"achDetails,omitempty"`
	CardDetails *ScheduleCardDetails `json:"cardDetails,omitempty"`
}

type ScheduleAchDetails struct {
	CompanyEntryDescription *string `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  *string `json:"originatingCompanyName,omitempty"`
}

type ScheduleCardDetails struct {
	DynamicDescriptor *string `json:"dynamicDescriptor,omitempty"`
}
