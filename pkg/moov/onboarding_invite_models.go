package moov

import "time"

// OnboardingInviteRequest is the payload for creating a new onboarding invite.
type OnboardingInviteRequest struct {
	ReturnURL         *string        `json:"returnURL,omitempty"`
	TermsOfServiceURL *string        `json:"termsOfServiceURL,omitempty"`
	Scopes            []string       `json:"scopes,omitempty"`
	GrantScopes       []string       `json:"grantScopes,omitempty"`
	Capabilities      []string       `json:"capabilities,omitempty"`
	FeePlanCodes      []string       `json:"feePlanCodes,omitempty"`
	Prefill           *CreateAccount `json:"prefill,omitempty"`
}

// OnboardingInvite represents an onboarding invite returned by the API.
type OnboardingInvite struct {
	Code              string         `json:"code"`
	Link              string         `json:"link"`
	ReturnURL         *string        `json:"returnURL,omitempty"`
	TermsOfServiceURL *string        `json:"termsOfServiceURL,omitempty"`
	Scopes            []string       `json:"scopes,omitempty"`
	GrantScopes       []string       `json:"grantScopes,omitempty"`
	Capabilities      []string       `json:"capabilities,omitempty"`
	FeePlanCodes      []string       `json:"feePlanCodes,omitempty"`
	RedeemedAccountID *string        `json:"redeemedAccountID,omitempty"`
	Prefill           *CreateAccount `json:"prefill,omitempty"`
	CreatedOn         time.Time      `json:"createdOn"`
	RevokedOn         *time.Time     `json:"revokedOn,omitempty"`
	RedeemedOn        *time.Time     `json:"redeemedOn,omitempty"`
}
