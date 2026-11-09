package testtools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const selfTestEnv = "MOOV_TESTTOOLS_SELFTEST"

func TestRandomForeignID(t *testing.T) {
	if os.Getenv(selfTestEnv) == "1" {
		fmt.Printf("RESULT:%s\n", foreignIDPrefix)
		return
	}

	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := RandomForeignID()
		require.NotEmpty(t, id)
		require.False(t, seen[id], "duplicate foreign ID %q", id)
		seen[id] = true
	}

	other := runSelf(t, "TestRandomForeignID", nil)
	require.NotEqual(t, foreignIDPrefix, other)
}

func TestSeededIdentifiers(t *testing.T) {
	if os.Getenv(selfTestEnv) == "1" {
		fmt.Printf("RESULT:%s|%s|%s|%s|%s\n", PARTNER_ID, PARTNER_WALLET_PM_ID, MERCHANT_ID, MERCHANT_WALLET_PM_ID, MERCHANT_WALLET_ID)
		return
	}

	defaults := []string{defaultPartnerID, defaultPartnerWalletPMID, defaultMerchantID, defaultMerchantWalletPMID, defaultMerchantWalletID}

	t.Run("unset keeps defaults", func(t *testing.T) {
		require.Equal(t, defaults, splitIDs(runSelf(t, "TestSeededIdentifiers", nil)))
	})

	t.Run("empty keeps defaults", func(t *testing.T) {
		overrides := map[string]string{
			"MOOV_PARTNER_ID":            "",
			"MOOV_PARTNER_WALLET_PM_ID":  "",
			"MOOV_MERCHANT_ID":           "",
			"MOOV_MERCHANT_WALLET_PM_ID": "",
			"MOOV_MERCHANT_WALLET_ID":    "",
		}
		require.Equal(t, defaults, splitIDs(runSelf(t, "TestSeededIdentifiers", overrides)))
	})

	t.Run("values override defaults", func(t *testing.T) {
		overrides := map[string]string{
			"MOOV_PARTNER_ID":            "env-partner-id",
			"MOOV_PARTNER_WALLET_PM_ID":  "env-partner-pm-id",
			"MOOV_MERCHANT_ID":           "env-merchant-id",
			"MOOV_MERCHANT_WALLET_PM_ID": "env-merchant-pm-id",
			"MOOV_MERCHANT_WALLET_ID":    "env-merchant-wallet-id",
		}
		want := []string{"env-partner-id", "env-partner-pm-id", "env-merchant-id", "env-merchant-pm-id", "env-merchant-wallet-id"}
		require.Equal(t, want, splitIDs(runSelf(t, "TestSeededIdentifiers", overrides)))
	})
}

func splitIDs(value string) []string {
	return strings.Split(value, "|")
}

// runSelf re-executes the test binary so the seeded identifiers resolve from
// the child process environment.
func runSelf(t *testing.T, testName string, overrides map[string]string) string {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$")
	cmd.Env = append(selfEnv(overrides), selfTestEnv+"=1")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "self subprocess failed: %s", out)

	for _, line := range strings.Split(string(out), "\n") {
		if value, ok := strings.CutPrefix(line, "RESULT:"); ok {
			return strings.TrimSpace(value)
		}
	}
	t.Fatalf("no RESULT line in subprocess output:\n%s", out)
	return ""
}

func selfEnv(overrides map[string]string) []string {
	strip := map[string]bool{selfTestEnv: true}
	for _, key := range []string{
		"MOOV_PARTNER_ID",
		"MOOV_PARTNER_WALLET_PM_ID",
		"MOOV_MERCHANT_ID",
		"MOOV_MERCHANT_WALLET_PM_ID",
		"MOOV_MERCHANT_WALLET_ID",
	} {
		strip[key] = true
	}

	env := make([]string, 0, len(os.Environ())+len(overrides))
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strip[key] {
			continue
		}
		env = append(env, entry)
	}
	for key, value := range overrides {
		env = append(env, key+"="+value)
	}
	return env
}
