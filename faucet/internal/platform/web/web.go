package web

import (
	"context"
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

func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
