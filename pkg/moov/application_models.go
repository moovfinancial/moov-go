package moov

import "time"

type ApplicationScope string

type Application struct {
	ApplicationID string `json:"applicationID,omitempty" otel:"application_id"`
	AccountID     string `json:"accountID,omitempty" otel:"account_id"`

	CreatedOn  time.Time  `json:"createdOn,omitempty" otel:"created_on"`
	UpdatedOn  time.Time  `json:"updatedOn,omitempty" otel:"updated_on"`
	DisabledOn *time.Time `json:"disabledOn,omitempty" otel:"disabled_on"`

	AllowedScopes []ApplicationScope `json:"allowedScopes,omitempty"`
	Name          string             `json:"name,omitempty"`
	Description   string             `json:"description,omitempty"`

	// don't use this, it's not used by the API
	AccountMode uint `json:"accountMode,omitempty"`
}

type CreateApplicationKey struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Origins     []string `json:"origins,omitempty"`
}

type ApplicationKeyWithSecret struct {
	ApplicationKeyID string `json:"applicationKeyID"`
	ApplicationID    string `json:"applicationID"`
	AccountID        string `json:"accountID"`

	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ClientId     string    `json:"clientId"`
	ClientSecret string    `json:"clientSecret"`
	LastUsed     time.Time `json:"lastUsed"`
	Origins      []string  `json:"origins"`

	CreatedOn time.Time `json:"createdOn"`

	UpdatedOn time.Time `json:"updatedOn"`

	DisabledOn *time.Time `json:"disabledOn"`
}
