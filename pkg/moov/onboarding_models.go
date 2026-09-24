package moov

import "time"

// This file contains the data models for the Onboarding Invites resource.
// See https://docs.moov.io/api/money-movement/onboarding/ for the API reference.

type OnboardingInvite struct {
	Code OnboardingInviteCode `json:"code"`
	// A unique URL, including the invite code, that the recipient can follow to redeem the invitation.
	Link string `json:"link"`
	// Optional URL to redirect the user to after they complete the onboarding process.
	ReturnURL string `json:"returnURL,omitempty"`
	// The terms of service URL set by the inviter.
	TermsOfServiceURL string `json:"termsOfServiceURL,omitempty"`
	// List of [scopes](https://docs.moov.io/api/authentication/scopes/) you request to use on this account. These values are used to determine what can be done with the account onboarded.
	Scopes []ApplicationScope `json:"scopes"`
	// List of [scopes](https://docs.moov.io/api/authentication/scopes/) you grant to allow being used by the new account on yourself. These values are used to determine what the account onboarded can do.
	GrantScopes []ApplicationScope `json:"grantScopes,omitempty"`
	// List of [capabilities](https://docs.moov.io/guides/accounts/capabilities/reference/) you intend to request for this account. These values are used to determine what information to collect from the user during onboarding.
	Capabilities []CapabilityName `json:"capabilities"`
	// List of fee plan codes to assign the account created by the invitee.
	FeePlanCodes []string `json:"feePlanCodes"`
	// The account ID of the account that redeemed the invite.
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
	// Optional URL to redirect the user to after they complete the onboarding process.
	ReturnURL string `json:"returnURL,omitempty"`
	// Optional URL to your organization's terms of service.
	TermsOfServiceURL string `json:"termsOfServiceURL,omitempty"`
	// List of [scopes](https://docs.moov.io/api/authentication/scopes/) you request to use on this account. These values are used to determine what can be done with the account onboarded.
	Scopes []ApplicationScope `json:"scopes"`
	// List of [scopes](https://docs.moov.io/api/authentication/scopes/) you grant to allow being used by the new account on yourself. These values are used to determine what the account onboarded can do.
	GrantScopes []ApplicationScope `json:"grantScopes,omitempty"`
	// List of [capabilities](https://docs.moov.io/guides/accounts/capabilities/reference/) you intend to request for this account. These values are used to determine what information to collect from the user during onboarding.
	Capabilities []CapabilityName `json:"capabilities"`
	// List of fee plan codes to assign the account created by the invitee.
	FeePlanCodes []string       `json:"feePlanCodes"`
	Prefill      *CreateAccount `json:"prefill,omitempty"`
}

// OnboardingPartnerAccount The account that created the onboarding invite.
type OnboardingPartnerAccount struct {
	// The account ID of the partner that created the invite.
	AccountID   string `json:"accountID"`
	AccountMode Mode   `json:"accountMode"`
	// The name of the Moov account used to create the onboarding invite.
	DisplayName string `json:"displayName"`
}
