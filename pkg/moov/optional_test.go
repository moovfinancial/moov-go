package moov

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptional_MarshalJSON(t *testing.T) {
	t.Run("set with value", func(t *testing.T) {
		opt := Set("hello")
		data, err := json.Marshal(opt)
		require.NoError(t, err)
		require.Equal(t, `"hello"`, string(data))
	})

	t.Run("set null", func(t *testing.T) {
		opt := SetNull[string]()
		data, err := json.Marshal(opt)
		require.NoError(t, err)
		require.Equal(t, `null`, string(data))
	})
}

func TestOptional_UnmarshalJSON(t *testing.T) {
	t.Run("value present", func(t *testing.T) {
		var opt Optional[string]
		err := json.Unmarshal([]byte(`"world"`), &opt)
		require.NoError(t, err)
		require.NotNil(t, opt.Value())
		require.Equal(t, "world", *opt.Value())
	})

	t.Run("null present", func(t *testing.T) {
		var opt Optional[string]
		err := json.Unmarshal([]byte(`null`), &opt)
		require.NoError(t, err)
		require.Nil(t, opt.Value())
	})
}

func TestOptional_Value(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var opt *Optional[int]
		require.Nil(t, opt.Value())
	})

	t.Run("Set", func(t *testing.T) {
		opt := Set(42)
		require.Equal(t, 42, *opt.Value())
	})

	t.Run("SetNull", func(t *testing.T) {
		opt := SetNull[int]()
		require.Nil(t, opt.Value())
	})
}

func TestOptional_IsNull(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var opt *Optional[int]
		require.False(t, opt.IsNull())
	})

	t.Run("Set value", func(t *testing.T) {
		opt := Set(42)
		require.False(t, opt.IsNull())
	})

	t.Run("SetNull", func(t *testing.T) {
		opt := SetNull[int]()
		require.True(t, opt.IsNull())
	})
}

func TestUpdateInvoice_MarshalJSON(t *testing.T) {
	t.Run("empty update produces empty object", func(t *testing.T) {
		u := UpdateInvoice{}
		data, err := json.Marshal(u)
		require.NoError(t, err)
		require.Equal(t, `{}`, string(data))
	})

	t.Run("set value is included", func(t *testing.T) {
		desc := "new description"
		u := UpdateInvoice{
			Description: &desc,
		}
		data, err := json.Marshal(u)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(data, &m))
		require.Equal(t, "new description", m["description"])
		require.NotContains(t, m, "dueDate")
		require.NotContains(t, m, "lineItems")
	})

	t.Run("null value is included as null", func(t *testing.T) {
		u := UpdateInvoice{
			DueDate: SetNull[time.Time](),
		}
		data, err := json.Marshal(u)
		require.NoError(t, err)
		require.Equal(t, `{"dueDate":null}`, string(data))
	})

	t.Run("mix of set, null, and unset", func(t *testing.T) {
		desc := "updated"
		u := UpdateInvoice{
			Description: &desc,
			DueDate:     SetNull[time.Time](),
			// Status, LineItems, InvoiceDate, TaxAmount left unset
		}
		data, err := json.Marshal(u)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(data, &m))
		require.Equal(t, "updated", m["description"])
		require.Contains(t, m, "dueDate")
		require.Nil(t, m["dueDate"])
		require.NotContains(t, m, "status")
		require.NotContains(t, m, "lineItems")
		require.NotContains(t, m, "invoiceDate")
		require.NotContains(t, m, "taxAmount")
	})
}
