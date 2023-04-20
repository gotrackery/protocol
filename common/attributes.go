package common

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/guregu/null.v4"
)

var (
	// ErrAttributesNoItem is returned when the attribute with given key is not present.
	ErrAttributesNoItem = errors.New("no item with given key")
	// ErrAttributesNotInt64 is returned when the attribute with given key is not of type int64.
	ErrAttributesNotInt64 = errors.New("value is not int64")
	// ErrAttributesNotFloat64 is returned when the attribute with given key is not of type float64.
	ErrAttributesNotFloat64 = errors.New("value is not float64")
	// ErrAttributesNotString is returned when the attribute with given key is not of type string.
	ErrAttributesNotString = errors.New("value is not string")
)

// Attributes is a generic attributes struct. It is used for key-value attributes for common purposes.
type Attributes map[string]interface{}

// GetType returns the type of the attribute.
func (a Attributes) GetType(key string) (string, error) {
	if v, ok := a[key]; ok {
		return reflect.TypeOf(v).String(), nil
	}
	// ToDo check valid types
	return "", ErrAttributesNoItem
}

// GetInt64 returns the int64 value of the attribute.
func (a Attributes) GetInt64(key string) (i int64, err error) {
	t, err := a.GetType(key)
	if err != nil {
		return 0, fmt.Errorf("failed to get type of attribute %s: %w", key, err)
	}
	if t != "int64" {
		return 0, ErrAttributesNotInt64
	}
	return a[key].(int64), nil //nolint:forcetypeassert
}

// GetFloat64 returns the float64 value of the attribute.
func (a Attributes) GetFloat64(key string) (f float64, err error) {
	t, err := a.GetType(key)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get type of attribute %s: %w", key, err)
	}
	if t != "float64" {
		return 0.0, ErrAttributesNotFloat64
	}
	return a[key].(float64), nil //nolint:forcetypeassert
}

// GetString returns the string value of the attribute.
func (a Attributes) GetString(key string) (f string, err error) {
	t, err := a.GetType(key)
	if err != nil {
		return "", fmt.Errorf("failed to get type of attribute %s: %w", key, err)
	}
	if t != "string" {
		return "", ErrAttributesNotString
	}
	return a[key].(string), nil //nolint:forcetypeassert
}

// AppendNullInt appends a null int to the attribute.
func (a Attributes) AppendNullInt(name string, i null.Int) Attributes {
	if !i.Valid {
		return a
	}

	m := a
	if m == nil {
		m = make(Attributes)
	}
	m[name] = i.Int64
	return m
}

// AppendNullFloat appends a null float to the attribute.
func (a Attributes) AppendNullFloat(name string, f null.Float) Attributes {
	if !f.Valid {
		return a
	}

	m := a
	if m == nil {
		m = make(Attributes)
	}
	m[name] = f.Float64
	return m
}

// AppendNullString appends a null string to the attribute.
func (a Attributes) AppendNullString(name string, s null.String) Attributes {
	if !s.Valid {
		return a
	}

	m := a
	if m == nil {
		m = make(Attributes)
	}
	m[name] = s.String
	return m
}

// AppendNullFloatSlice appends a null float slice to the attribute.
func (a Attributes) AppendNullFloatSlice(name string, s []null.Float) Attributes {
	m := a
	if m == nil {
		m = make(Attributes)
	}
	for i, v := range s {
		if v.Valid {
			m[fmt.Sprintf("%s_%d", name, i)] = v.Float64
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}
