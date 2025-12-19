package moov

// This allows for a connections to a partner account to be shared with another account on that partner account.
//
// Allows for going from the customer and merchant only being connected to the partner to being able for the merchant to be connected to the customer as well.
//   customer ---> partner <---- merchant
//
// to
//   +--> customer ----> partner <---- merchant --+
//   +--------------------------------------------+
//

type ShareConnectionRequest struct {
	// This is the account to be accessing the subject account who is defined in the path.
	PrincipalAccountID string `json:"principalAccountID"`

	// These are the scopes that the principal account is allowed to access on the subject account.
	AllowScopes []string `json:"allowScopes"`
}

type ShareConnectionResponse struct {
	PrincipalAccountID string `json:"principalAccountID"`
	SubjectAccountID   string `json:"subjectAccountID"`

	// List of scopes that were allowed to be shared. It'll be a the list of all scopes or a subset of the allowed scopes specified in the request.
	Scopes []string `json:"scopes"`
}
