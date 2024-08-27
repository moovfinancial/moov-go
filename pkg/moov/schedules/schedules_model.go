package schedules

import "time"

type Schedule struct {
	// prod or sandbox
	Mode string `json:"mode,omitempty" spanner:"mode"`

	ScheduleID string `json:"scheduleID,omitempty" spanner:"schedule_id"`

	// This is the account ID of the source transfers that the transfer will run using.
	SourceAccountID string `json:"sourceAccountID,omitempty" spanner:"source_account_id"`

	// This is the destination ID of the destination transfers that the transfer will run using.
	DestinationAccountID string `json:"destinationAccountID,omitempty" spanner:"destination_account_id"`

	// This is the partner ID that the transfer will run using.
	PartnerAccountID string `json:"partnerAccountID,omitempty" spanner:"partner_account_id"`

	// AccountID of the account that created it and is allowed to update it.
	OwnerID string `json:"ownerID,omitempty" spanner:"owner_id"`

	// Description of what this schedule is
	Description string `json:"description,omitempty" spanner:"description"`

	// If specified will generate Scheduled transfers based on its configuration
	RecurTransfer *RecurTransfer `json:"recurTransfer,omitempty" spanner:"recur_transfer"`

	// List of all generated and manually added transfers to be made.
	Occurrences []TransferOccurrence `json:"occurrences,omitempty" spanner:"-"`

	CreatedOn  time.Time  `json:"createdOn,omitempty" spanner:"created_on"`
	UpdatedOn  time.Time  `json:"updatedOn,omitempty" spanner:"updated_on"`
	DisabledOn *time.Time `json:"disabledOn,omitempty"  spanner:"disabled_on"`
}

func (s Schedule) ToUpsertSchedule() UpsertSchedule {
	occs := make([]UpsertTransferOccurrence, len(s.Occurrences))
	for i, occ := range s.Occurrences {
		occs[i] = UpsertTransferOccurrence{
			OccurrenceID: occ.OccurrenceID,
			Transfer:     occ.Transfer,
			RunOn:        occ.RunOn,
			Cancelled:    nil,
		}
	}

	return UpsertSchedule{
		Description:   s.Description,
		RecurTransfer: s.RecurTransfer,
		Occurrences:   occs,
	}
}

// https://www.rfc-editor.org/rfc/rfc5545#section-3.3.10
type RecurTransfer struct {
	// Transfer values to use to create the transfer based on the recurRule
	// When changed, should just modify the transfer of the schedules
	Transfer Transfer `json:"transfer,omitempty"`

	// https://www.rfc-editor.org/rfc/rfc5545#section-3.3.10
	RecurrenceRule string `json:"recurrenceRule,omitempty"`

	// If the recurrence rule ends up being indefinite
	Indefinite bool `json:"indefinite,omitempty"`
}

type TransferOccurrence struct {
	ScheduleID string `json:"scheduleID,omitempty" spanner:"schedule_id"`

	// Unique ID for updating a specific occurrence
	OccurrenceID string `json:"occurrenceID,omitempty" spanner:"occurrence_id"`

	// Mode to run the occurence under
	Mode string `json:"mode,omitempty" spanner:"mode"`

	// Transfer details that will be used.
	Transfer Transfer `json:"transfer,omitempty" spanner:"transfer"`

	// If this scheduled transfer was generated or manually added for say a correction
	// If a new interval is specified, all un-ran generated transfers will be re-generated
	Generated bool `json:"generated,omitempty" spanner:"generated"`

	// Flag if this is part of an indefinite schedule
	Indefinite bool `json:"indefinite,omitempty" spanner:"indefinite"`

	// Modified since generated. This could be switching just a single payment method
	Modified bool `json:"modified,omitempty" spanner:"modified"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty" spanner:"run_on"`

	// When the transfer was kicked off. If nil, hasn't ran. Normalize to UTC.
	RanOn *time.Time `json:"ranOn,omitempty" spanner:"ran_on"`

	// Ability to cancel this specific transfer from running
	CancelledOn *time.Time `json:"cancelledOn,omitempty" spanner:"cancelled_on"`

	// ID of the transfer that ran
	TransferID     *string `json:"transferID,omitempty" spanner:"transfer_id"`
	TransferStatus *string `json:"transferStatus,omitempty" spanner:"transfer_status"`
}

type UpsertSchedule struct {
	// Description of what this schedule is
	Description string `json:"description,omitempty"`

	// If specified will generate Scheduled transfers based on its configuration
	RecurTransfer *RecurTransfer `json:"recurTransfer,omitempty"`

	// On creating the schedule we can use these occurrences as they planned the schedule
	Occurrences []UpsertTransferOccurrence `json:"occurrences,omitempty"`
}

type UpsertTransferOccurrence struct {
	// Leave empty to add a new occurrence or set to the ID of the occurrence to change.
	OccurrenceID string `json:"occurrenceID,omitempty" spanner:"occurrence_id"`

	// Transfer details that will be used.
	Transfer Transfer `json:"transfer,omitempty"`

	// Time to kick off the run. Normalize to UTC.
	RunOn time.Time `json:"runOn,omitempty"`

	// If nil, cancelledOn will be unchanged. If set true, it will be cancelled. If set false and hasn't ran yet will be uncancelled
	Cancelled *bool `json:"cancelled,omitempty"`
}
