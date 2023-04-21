package dbtp

import (
	"encoding/json"
	"time"
)

// JsonTime basic time value to json format
type JsonTime time.Time

var JsonTimeFormat = "2006-01-02 15:04:05"

// MarshalJSON
func (jsonTime JsonTime) MarshalJSON() ([]byte, error) {
	simple := time.Time(jsonTime).Format(JsonTimeFormat)
	return json.Marshal(simple)
}

// UnmarshalJSON
func (jsonTime *JsonTime) UnmarshalJSON(bt []byte) error {
	var val string
	if err := json.Unmarshal(bt, &val); err != nil {
		return err
	}
	t, err := time.Parse(JsonTimeFormat, val)
	if err != nil {
		return err
	}
	*jsonTime = JsonTime(t)
	return nil
}

func NowJsonTime() *JsonTime {
	j := JsonTime(time.Now())
	return &j
}
