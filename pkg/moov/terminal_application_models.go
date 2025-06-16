package moov

type TerminalApplicationStatus string

// List of TerminalApplicationStatus
const (
	TerminalApplicationStatus_Pending  TerminalApplicationStatus = "pending"
	TerminalApplicationStatus_Enabled  TerminalApplicationStatus = "enabled"
	TerminalApplicationStatus_Disabled TerminalApplicationStatus = "disabled"
)

type TerminalApplicationPlatform string

// List of TerminalApplicationPlatform
const (
	TerminalApplicationPlatform_Android TerminalApplicationPlatform = "android"
	TerminalApplicationPlatform_iOS     TerminalApplicationPlatform = "ios"
)

type TerminalApplication struct {
	// UUID
	TerminalApplicationID string                      `json:"terminalApplicationID"`
	Status                TerminalApplicationStatus   `json:"status"`
	Platform              TerminalApplicationPlatform `json:"platform"`
	AppBundleID           string                      `json:"appBundleID,omitempty"`
	PackageName           string                      `json:"packageName,omitempty"`
	Sha256Digest          string                      `json:"sha256Digest,omitempty"`
	VersionCode           string                      `json:"versionCode,omitempty"`
}

type TerminalApplicationRequest struct {
	Platform TerminalApplicationPlatform `json:"platform"`
	// The app bundle identifier of the terminal application. Required if platform is `ios`.
	AppBundleID string `json:"appBundleID,omitempty"`
	// The app package name of the terminal application. Required if platform is `android`.
	PackageName string `json:"packageName,omitempty"`
	// The app version of the terminal application. Required if platform is `android`.
	Sha256Digest string `json:"sha256Digest,omitempty"`
	// The app version of the terminal application. Required if platform is `android`.
	VersionCode string `json:"versionCode,omitempty"`
}

type TerminalApplicationVersion struct {
	Version string `json:"version"`
}
