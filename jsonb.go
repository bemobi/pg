package pg

import (
	"encoding/json"
	"errors"

	"github.com/guregu/null"
)

type JSONB null.String

func (j *JSONB) Encode(src interface{}) error {
	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}
	*j = JSONB(null.NewString(string(jsonData), true))
	return nil
}

func (j *JSONB) Decode(dst interface{}) error {
	if !j.NullString.Valid {
		return errors.New("Empty JSON data")
	}
	return json.Unmarshal([]byte(j.NullString.String), dst)
}

// MarshalJSON implements the `json.Marshaller` interface
func (j *JSONB) MarshalJSON() ([]byte, error) {
	if j.NullString.Valid {
		return []byte(j.NullString.String), nil
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the `json.Unmarshaler` interface
func (j *JSONB) UnmarshalJSON(bytes []byte) error {
	if len(bytes) > 0 {
		j.NullString.String = string(bytes)
		j.NullString.Valid = true
	} else {
		j.NullString.String = ""
		j.NullString.Valid = false
	}
	return nil
}
