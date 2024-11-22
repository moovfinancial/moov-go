package moov

import "time"

// Representative Describes a business representative.
type Representative struct {
	RepresentativeID string `json:"representativeID,omitempty"`
	Name             Name   `json:"name,omitempty"`
	Phone            *Phone `json:"phone,omitempty"`
	// Email address.
	Email   string   `json:"email,omitempty"`
	Address *Address `json:"address,omitempty"`
	// Indicates whether this representative's birth date has been provided.
	BirthDateProvided bool `json:"birthDateProvided,omitempty"`
	// Indicates whether a government ID (SSN, ITIN, etc.) has been provided for this representative.
	GovernmentIDProvided bool              `json:"governmentIDProvided,omitempty"`
	Responsibilities     *Responsibilities `json:"responsibilities,omitempty"`
	CreatedOn            time.Time         `json:"createdOn,omitempty"`
	UpdatedOn            time.Time         `json:"updatedOn,omitempty"`
	DisabledOn           *time.Time        `json:"disabledOn,omitempty"`
}

type CreateRepresentative struct {
	Name             Name              `json:"name"`
	Phone            *Phone            `json:"phone,omitempty"`
	Email            string            `json:"email,omitempty"`
	Address          *Address          `json:"address,omitempty"`
	BirthDate        *Date             `json:"birthDate,omitempty"`
	GovernmentID     *GovernmentID     `json:"governmentID,omitempty"`
	Responsibilities *Responsibilities `json:"responsibilities,omitempty"`
}

type UpdateRepresentative struct {
	Name             Name              `json:"name"`
	Phone            *Phone            `json:"phone,omitempty"`
	Email            string            `json:"email,omitempty"`
	Address          *Address          `json:"address,omitempty"`
	BirthDate        *Date             `json:"birthDate,omitempty"`
	GovernmentID     *GovernmentID     `json:"governmentID,omitempty"`
	Responsibilities *Responsibilities `json:"responsibilities,omitempty"`
}
