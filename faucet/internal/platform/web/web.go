package web

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}
	return nil
}
