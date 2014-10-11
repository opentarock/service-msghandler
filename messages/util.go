package messages

import "encoding/json"

func Marshal(v interface{}) string {
	r, err := json.Marshal(v)
	if err != nil {
		panic("Unable to marshal message")
	}
	return string(r)
}
