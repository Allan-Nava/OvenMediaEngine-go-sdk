package ovenmedia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Time is a timestamp from OME. It accepts both RFC 3339 offsets
// ("+09:00") and the colon-less form some OME versions print ("+0900");
// an empty string or null decodes to the zero time.
type Time struct {
	time.Time
}

var timeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999-0700",
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		t.Time = time.Time{}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("ovenmedia: timestamp %s: %w", b, err)
	}
	if s == "" {
		t.Time = time.Time{}
		return nil
	}
	for _, layout := range timeLayouts {
		if parsed, err := time.Parse(layout, s); err == nil {
			t.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("ovenmedia: unrecognised timestamp %q", s)
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(t.Format(time.RFC3339Nano))
}

// FlexInt64 is a number that OME sends either as a JSON number or, in older
// versions, as a string ("2500000"). Bitrates are the usual case.
type FlexInt64 int64

func (n *FlexInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		*n = 0
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*n = 0
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("ovenmedia: number %q: %w", s, err)
		}
		*n = FlexInt64(v)
		return nil
	}
	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*n = FlexInt64(v)
	return nil
}
