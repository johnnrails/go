package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type NullInt64 struct{ sql.NullInt64 }

func (ni *NullInt64) Marshal() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	jsn, err := json.Marshal(ni.Int64)
	if err != nil {
		return jsn, fmt.Errorf("MySQL Could not marshal NullInt64: %w", err)
	}
	return jsn, nil
}

func (ni *NullInt64) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, &ni.Int64); err != nil {
		return fmt.Errorf("MySQL NullInt64 unmarshal error: %w", err)
	}
	ni.Valid = true
	return nil
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct{ sql.NullBool }

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}

	jsn, err := json.Marshal(nb.Bool)
	if err != nil {
		return jsn, fmt.Errorf("MySQL could not marshal NullBool: %w", err)
	}

	return jsn, nil
}

// UnmarshalJSON for NullBool
func (nb NullBool) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &nb.Bool); err != nil {
		return fmt.Errorf("MySQL NullBool unmarshal error: %w", err)
	}

	nb.Valid = true

	return nil
}

type NullFloat64 struct{ sql.NullFloat64 }

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}

	jsn, err := json.Marshal(nf.Float64)
	if err != nil {
		return jsn, fmt.Errorf("MySQL could not marshal NullFloat64: %w", err)
	}

	return jsn, nil
}

func (nf NullFloat64) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &nf.Float64); err != nil {
		return fmt.Errorf("MySQL NullFloat64 unmarshal error: %w", err)
	}
	nf.Valid = true
	return nil
}

type NullString struct{ sql.NullString }

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	jsn, err := json.Marshal(ns.String)
	if err != nil {
		return jsn, fmt.Errorf("MySQL could not marshal NullString: %w", err)
	}

	return jsn, nil
}

func (ns NullString) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &ns.String); err != nil {
		return fmt.Errorf("MySQL NullString unmarshal error: %w", err)
	}

	ns.Valid = true

	return nil
}

type NullTime struct{ sql.NullTime }

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (nt NullTime) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &nt.Time); err != nil {
		return fmt.Errorf("MySQL NullTime unmarshal error: %w", err)
	}
	nt.Valid = true
	return nil
}
