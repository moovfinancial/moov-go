package moov_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestPercentageRate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  moov.PercentageRate
	}{
		{"json number", `{"percentageRate":0.0195}`, "0.0195"},
		{"trailing zeros preserved", `{"percentageRate":0.019500000}`, "0.019500000"},
		{"nine decimal places", `{"percentageRate":1.123456789}`, "1.123456789"},
		{"integer", `{"percentageRate":2}`, "2"},
		{"negative", `{"percentageRate":-0.0195}`, "-0.0195"},
		{"exponent", `{"percentageRate":1.95e-2}`, "1.95e-2"},
		{"zero", `{"percentageRate":0}`, "0"},
		{"field absent", `{"programName":"Visa CPS"}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fee moov.InterchangeProgramFee
			require.NoError(t, json.Unmarshal([]byte(tt.input), &fee))
			require.Equal(t, tt.want, fee.PercentageRate)
		})
	}
}

func TestPercentageRate_UnmarshalJSON_Errors(t *testing.T) {
	t.Run("bool is a type error", func(t *testing.T) {
		var fee moov.InterchangeProgramFee
		err := json.Unmarshal([]byte(`{"percentageRate":true}`), &fee)
		require.Error(t, err)

		var typeErr *json.UnmarshalTypeError
		require.ErrorAs(t, err, &typeErr)
	})

	t.Run("empty string is not a valid number", func(t *testing.T) {
		var fee moov.InterchangeProgramFee
		err := json.Unmarshal([]byte(`{"percentageRate":""}`), &fee)
		require.ErrorContains(t, err, "invalid number literal")
	})

	t.Run("non-numeric string is not a valid number", func(t *testing.T) {
		var fee moov.InterchangeProgramFee
		err := json.Unmarshal([]byte(`{"percentageRate":"one percent"}`), &fee)
		require.ErrorContains(t, err, "invalid number literal")
	})
}

func TestPercentageRate_MarshalJSON(t *testing.T) {
	t.Run("encodes as an unquoted number", func(t *testing.T) {
		data, err := json.Marshal(moov.PercentageRate("0.0195"))
		require.NoError(t, err)
		require.Equal(t, `0.0195`, string(data))
	})

	t.Run("precision is not rounded through a float", func(t *testing.T) {
		data, err := json.Marshal(moov.PercentageRate("1.123456789"))
		require.NoError(t, err)
		require.Equal(t, `1.123456789`, string(data))
	})

	t.Run("struct field is unquoted", func(t *testing.T) {
		data, err := json.Marshal(moov.InterchangeProgramFee{PercentageRate: "0.0195"})
		require.NoError(t, err)
		require.Contains(t, string(data), `"percentageRate":0.0195`)
	})

	t.Run("zero value is omitted", func(t *testing.T) {
		data, err := json.Marshal(moov.InterchangeProgramFee{ProgramName: "Visa CPS"})
		require.NoError(t, err)
		require.NotContains(t, string(data), "percentageRate")
	})

	t.Run("round trip is byte identical", func(t *testing.T) {
		var fee moov.InterchangeProgramFee
		require.NoError(t, json.Unmarshal([]byte(`{"percentageRate":0.019500000}`), &fee))

		data, err := json.Marshal(fee.PercentageRate)
		require.NoError(t, err)
		require.Equal(t, `0.019500000`, string(data))
	})
}

func TestPercentageRate_Float64(t *testing.T) {
	var fee moov.InterchangeProgramFee
	require.NoError(t, json.Unmarshal([]byte(`{"percentageRate":0.0195}`), &fee))

	f, err := json.Number(fee.PercentageRate).Float64()
	require.NoError(t, err)
	require.InDelta(t, 0.0195, f, 1e-9)
}

// statementWithInterchangePrograms mirrors a statements response carrying a populated interchangePrograms array.
// The array is omitted when empty, which is why the field's original string/number mismatch went unnoticed.
const statementWithInterchangePrograms = `{
  "statementID": "8e5b8c4f-6a1d-4f7e-9b3a-2c1d4e5f6a7b",
  "statementName": "July 2026",
  "billingPeriodStartDateTime": "2026-07-01T00:00:00Z",
  "billingPeriodEndDateTime": "2026-08-01T00:00:00Z",
  "summary": {
    "cardAcquiring": {
      "feeAmount": { "currency": "USD", "valueDecimal": "1204.567890000" },
      "interchangeFees": {
        "visa": { "currency": "USD", "valueDecimal": "900.120000000" },
        "mastercard": { "currency": "USD", "valueDecimal": "304.447890000" },
        "discover": { "currency": "USD", "valueDecimal": "0.000000000" },
        "americanExpress": { "currency": "USD", "valueDecimal": "0.000000000" }
      }
    },
    "total": { "currency": "USD", "valueDecimal": "1204.567890000" }
  },
  "cardAcquiringFees": {
    "visa": {
      "interchangePrograms": [
        {
          "count": 1250,
          "perItemRate": { "currency": "USD", "valueDecimal": "0.100000000" },
          "percentageRate": 0.0195,
          "programName": "Visa CPS Retail",
          "transferVolume": { "currency": "USD", "valueDecimal": "45000.000000000" },
          "total": { "currency": "USD", "valueDecimal": "1002.500000000" }
        },
        {
          "count": 40,
          "perItemRate": { "currency": "USD", "valueDecimal": "0.050000000" },
          "percentageRate": 1.123456789,
          "programName": "Visa CPS Card Not Present",
          "transferVolume": { "currency": "USD", "valueDecimal": "1200.000000000" },
          "total": { "currency": "USD", "valueDecimal": "15.481481468" }
        }
      ],
      "total": { "amount": { "currency": "USD", "valueDecimal": "1017.981481468" }, "count": 1290 }
    }
  }
}`

func TestStatement_InterchangePrograms(t *testing.T) {
	var statement moov.Statement
	err := strictDecoder(strings.NewReader(statementWithInterchangePrograms), "application/json", &statement)
	require.NoError(t, err)

	require.NotNil(t, statement.CardAcquiringFees)
	programs := statement.CardAcquiringFees.Visa.InterchangePrograms
	require.NotNil(t, programs)
	require.Len(t, *programs, 2)

	require.Equal(t, moov.PercentageRate("0.0195"), (*programs)[0].PercentageRate)
	require.Equal(t, "Visa CPS Retail", (*programs)[0].ProgramName)

	require.Equal(t, moov.PercentageRate("1.123456789"), (*programs)[1].PercentageRate)
	require.Equal(t, int64(40), (*programs)[1].Count)
}
