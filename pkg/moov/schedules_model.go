package moov

import "time"

type Schedule struct {
	// prod or sandbox
	Mode string `json:"mode,omitempty"`

	// Unique ID of the schedule
	ScheduleID string `json:"scheduleID,omitempty"`

	// This is the account ID of the source transfers that the transfer will run using.
	SourceAccountID string `json:"sourceAccountID,omitempty"`

	// This is the destination ID of the destination transfers that the transfer will run using.
	DestinationAccountID string `json:"destinationAccountID,omitempty"`

	// This is the partner ID that the transfer will run using.
	PartnerAccountID string `json:"partnerAccountID,omitempty"`

	// AccountID of the account that created it and is allowed to update it.
	OwnerAccountID string `json:"ownerAccountID,omitempty"`

	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	Recur *Recur `json:"recur,omitempty"`

	// List of all generated and manually added transfers to be made.
	Occurrences []Occurrence `json:"occurrences,omitempty"`

	// Date created
	CreatedOn time.Time `json:"createdOn,omitempty"`

	// Date it was last updated for any reason
	UpdatedOn time.Time `json:"updatedOn,omitempty"`

	// When schedule has been disabled and all occurrences canceled
	DisabledOn *time.Time `json:"disabledOn,omitempty"`
}

func (s Schedule) ToUpdateSchedule() UpdateSchedule {
	upsOccs := make([]UpdateOccurrence, len(s.Occurrences))
	for i, occ := range s.Occurrences {
		upsOccs[i] = UpdateOccurrence{
			OccurrenceID: &occ.OccurrenceID,
			RunTransfer:  occ.RunTransfer,
			RunOn:        occ.RunOn,
			Canceled:     nil,
		}
	}

	return UpdateSchedule{
		Description: s.Description,
		Recur:       s.Recur,
		Occurrences: upsOccs,
	}
}

// https://www.rfc-editor.org/rfc/rfc5545#section-3.3.10
type Recur struct {
	// If omitted the start time for the occurrence will be the timestamp of when the schedule was created.
	Start *time.Time `json:"start,omitempty"`

	// This is the recurrence rule that is used to generate occurrences.
	// Generator available here: https://jkbrzt.github.io/rrule/
	// You can read the details of the format here: https://www.rfc-editor.org/rfc/rfc5545#section-3.3.10
	RecurrenceRule string `json:"recurrenceRule,omitempty"`

	// RunTransfer values to use to create the transfer based on the recurRule
	// When changed, should just modify the transfer of the schedules
	RunTransfer RunTransfer `json:"runTransfer,omitempty"`

	// If the recurrence rule ends up being indefinite
	Indefinite bool `json:"indefinite,omitempty"`
}

type Occurrence struct {
	ScheduleID string `json:"scheduleID,omitempty"`

	// Unique ID for updating a specific occurrence
	OccurrenceID string `json:"occurrenceID,omitempty"`

	// Mode to run the occurrence under
	Mode string `json:"mode,omitempty"`

	// If this scheduled transfer was generated or manually added for say a correction
	// If a new interval is specified, all un-ran generated transfers will be re-generated
	Generated bool `json:"generated,omitempty"`

	// Flag if this is part of an indefinite schedule
	Indefinite bool `json:"indefinite,omitempty"`

	// Modified since generated. This could be switching just a single payment method
	Modified bool `json:"modified,omitempty"`

	// Ability to cancel this specific transfer from running
	CanceledOn *time.Time `json:"canceledOn,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`

	// RunTransfer details that will be used.
	RunTransfer RunTransfer `json:"runTransfer,omitempty"`

	// When the transfer was kicked off. If nil, hasn't ran. Normalize to UTC.
	RanOn *time.Time `json:"ranOn,omitempty"`

	// ID of the transfer that ran
	RunTransferID *string `json:"ranTransferID,omitempty"`

	// Status of the running occurrence
	Status *string `json:"status,omitempty"`

	// Descriptive message of why it errored.
	Error *OccurrenceError `json:"error,omitempty" spanner:"error" otel:"error"`
}

// OccurrenceError is where we log any errors or failures that could happen from running the occurrence.
type OccurrenceError struct {
	Message string `json:"message,omitempty" otel:"message"`
}

type CreateSchedule struct {
	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	Recur *Recur `json:"recur,omitempty"`

	// On creating the schedule we can use these occurrences as they planned the schedule
	Occurrences []CreateOccurrence `json:"occurrences,omitempty"`
}

type CreateOccurrence struct {
	// RunTransfer details that will be used.
	RunTransfer RunTransfer `json:"runTransfer,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`
}

type UpdateSchedule struct {
	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	Recur *Recur `json:"recur,omitempty"`

	// On creating the schedule we can use these occurrences as they planned the schedule
	Occurrences []UpdateOccurrence `json:"occurrences,omitempty"`
}

type UpdateOccurrence struct {
	// Leave empty to add a new occurrence or set to the ID of the occurrence to change.
	OccurrenceID *string `json:"occurrenceID,omitempty"`

	// RunTransfer details that will be used.
	RunTransfer RunTransfer `json:"runTransfer,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`

	// If nil, canceledOn will be unchanged. If set true, it will be canceled. If set false and hasn't ran yet it will be resumed
	Canceled *bool `json:"canceled,omitempty"`
}

type RunTransfer struct {
	Description string         `json:"description,omitempty"`
	Amount      ScheduleAmount `json:"amount,omitempty"`

	PartnerAccountID string                `json:"partnerAccountID,omitempty"`
	Source           SchedulePaymentMethod `json:"source,omitempty"`
	Destination      SchedulePaymentMethod `json:"destination,omitempty"`
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
