package moov

import "encoding/json"

// Optional represents a field that can be in one of three states:
//   - Not set: a nil *Optional is omitted from JSON via omitempty
//   - Explicitly null: the field is serialized as null in JSON (use SetNull)
//   - Set to a value: the field is serialized as the value in JSON (use Set)
//
// This is useful for PATCH endpoints where nil/*T cannot distinguish between
// "don't update this field" and "unset this field".
//
// Use *Optional[T] with `json:",omitempty"` on struct fields so that nil
// pointers (not provided) are omitted automatically.
type Optional[T any] struct {
	value *T
}

// Set creates an Optional with a value. The field will be serialized to its JSON value.
func Set[T any](v T) *Optional[T] {
	return &Optional[T]{value: &v}
}

// SetNull creates an Optional that is explicitly null.
// The field will be serialized as null in JSON.
func SetNull[T any]() *Optional[T] {
	return &Optional[T]{}
}

// Value returns a pointer to the value, or nil if null.
func (o *Optional[T]) Value() *T {
	if o == nil {
		return nil
	}
	return o.value
}

// IsNull returns true if the field is explicitly set to null.
func (o *Optional[T]) IsNull() bool {
	return o != nil && o.value == nil
}

// MarshalJSON implements json.Marshaler.
// If set with a value, marshals the value. If set to null, marshals as null.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*o.value)
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.value = nil
		return nil
	}
	o.value = new(T)
	return json.Unmarshal(data, o.value)
}
