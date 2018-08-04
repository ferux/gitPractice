package hidstruct

import (
	"encoding/json"
)

var (
	//FilterMessage stores replace string
	FilterMessage = "[FILTERED]"
)

// StringSafe ensures data will be hide while displaying
type StringSafe string

func (s StringSafe) String() string {
	return FilterMessage
}

// GoString for %#v display
func (s StringSafe) GoString() string {
	return FilterMessage
}

// MarshalJSON for marshaing data
func (s StringSafe) MarshalJSON() ([]byte, error) {
	return json.Marshal(FilterMessage)
}
