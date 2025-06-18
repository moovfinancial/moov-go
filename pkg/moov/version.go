package moov

import (
	"fmt"
)

const VersionHeader = "X-Moov-Version"

var (
	// Only selectable if beta is set.
	WorkInProgress = Version{"v0001.00.00"}

	// Pre-versioned API
	PreVersioning = Version{"v2024.00.00"}

	Version2025_01 = NewVersion(2025, 1, 0)
	Version2025_04 = NewVersion(2025, 4, 0)
	Version2025_07 = NewVersion(2025, 7, 0)
	Version2025_10 = NewVersion(2025, 10, 0)

	Version2026_01 = NewVersion(2026, 1, 0)
	Version2026_04 = NewVersion(2026, 4, 0)
	Version2026_07 = NewVersion(2026, 7, 0)
	Version2026_10 = NewVersion(2026, 10, 0)

	// Selects the latest version that isn't Beta
	Latest = Version{"v9000.00.00"}
)

type Version struct {
	version string
}

func (t Version) String() string {
	return t.version
}

// Best to set to the anticipated release date or into the future until release.
func NewVersion(year int, month int, build int) Version {
	return Version{fmt.Sprintf("v%04d.%02d.%02d", year, month, build)}
}
