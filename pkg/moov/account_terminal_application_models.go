package moov

type AccountTerminalApplication struct {
	AccountID             string `json:"accountID"`
	TerminalApplicationID string `json:"terminalApplicationID"`
}

type LinkAccountTerminalApplicationRequest struct {
	TerminalApplicationID string `json:"terminalApplicationID"`
}

type AccountTerminalApplicationConfiguration struct {
	Configuration string `json:"configuration"`
}
