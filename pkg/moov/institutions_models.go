package moov

type InstitutionsSearchResponse struct {
	Ach  []ACHInstitution  `json:"ach"`
	Rtp  []RTPInstitution  `json:"rtp"`
	Wire []WireInstitution `json:"wire"`
}

type ACHInstitution struct {
	Name          string   `json:"name"`
	RoutingNumber string   `json:"routingNumber"`
	Address       *Address `json:"address,omitempty"`
	Contact       *Contact `json:"contact,omitempty"`
}

type RTPInstitution struct {
	Name          string      `json:"name"`
	RoutingNumber string      `json:"routingNumber"`
	Services      RTPServices `json:"services"`
}

type RTPServices struct {
	// Can the institution receive payments
	ReceivePayments bool `json:"receivePayments"`
	// Can the institution receive request for payment messages
	ReceiveRequestForPayment bool `json:"receiveRequestForPayment"`
}

type WireInstitution struct {
	Name          string       `json:"name"`
	RoutingNumber string       `json:"routingNumber"`
	Address       *Address     `json:"address,omitempty"`
	Services      WireServices `json:"services"`
}

type WireServices struct {
	// The institution's capability to process standard Fedwire funds transfers.
	FundsTransferStatus bool `json:"fundsTransferStatus"`
	// The institution's capability for settlement-only transfers.
	FundsSettlementOnlyStatus bool `json:"fundsSettlementOnlyStatus"`
	// The institution's capability to handle transfers of securities.
	BookEntrySecuritiesTransferStatus bool `json:"bookEntrySecuritiesTransferStatus"`
}
