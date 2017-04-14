package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ur/errors"
)

type jsonError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (j jsonError) Write(w http.ResponseWriter) {

	if j.StatusCode > 0 {
		w.WriteHeader(j.StatusCode)
	}

	if jsn, jsonErr := json.Marshal(j); jsonErr == nil {
		w.Write(jsn)
	} else {
		log.Printf("api.JSONError: Failed to serialise JSON: %s", jsonErr)
	}
}

func CodedJSONErrorMessage(w http.ResponseWriter, statusCode int, desc string, args ...interface{}) (err error) {

	err = fmt.Errorf(desc, args...)

	errObject := jsonError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}

	errObject.Write(w)

	return
}

func JSONError(w http.ResponseWriter, err error) error {

	errObject := jsonError{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}

	switch err.(type) {
	case errors.User:
		errObject.StatusCode = http.StatusBadRequest
	}

	errObject.Write(w)

	return err
}
