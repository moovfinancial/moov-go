package moov

import "time"

type OnboardingInvite struct {
	Code              OnboardingInviteCode      `json:"code"`
	Link              string                    `json:"link"`
	ReturnURL         string                    `json:"returnURL,omitempty"`
	TermsOfServiceURL string                    `json:"termsOfServiceURL,omitempty"`
	Scopes            []ApplicationScope        `json:"scopes"`
	GrantScopes       []ApplicationScope        `json:"grantScopes,omitempty"`
	Capabilities      []CapabilityName          `json:"capabilities"`
	FeePlanCodes      []string                  `json:"feePlanCodes"`
	RedeemedAccountID string                    `json:"redeemedAccountID,omitempty"`
	Prefill           *CreateAccount            `json:"prefill,omitempty"`
	Partner           *OnboardingPartnerAccount `json:"partner,omitempty"`
	CreatedOn         time.Time                 `json:"createdOn"`
	RevokedOn         *time.Time                `json:"revokedOn,omitempty"`
	RedeemedOn        *time.Time                `json:"redeemedOn,omitempty"`
}

// OnboardingInviteCode A unique code that identifies an onboarding invite.
type OnboardingInviteCode string

// OnboardingInviteRequest Request to create an onboarding invite.
type OnboardingInviteRequest struct {
	ReturnURL         string             `json:"returnURL,omitempty"`
	TermsOfServiceURL string             `json:"termsOfServiceURL,omitempty"`
	Scopes            []ApplicationScope `json:"scopes"`
	GrantScopes       []ApplicationScope `json:"grantScopes,omitempty"`
	Capabilities      []CapabilityName   `json:"capabilities"`
	FeePlanCodes      []string           `json:"feePlanCodes"`
	Prefill           *CreateAccount     `json:"prefill,omitempty"`
}

// OnboardingPartnerAccount The account that created the onboarding invite.
type OnboardingPartnerAccount struct {
	AccountID   string `json:"accountID"`
	AccountMode Mode   `json:"accountMode"`
	DisplayName string `json:"displayName"`
}
