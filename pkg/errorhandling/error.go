package errorhandling

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// apierror represents the response body in case of error
type apierror struct {
	// Code is a short string indicating the error code
	code string
	// Message is a human readable message providing more details about the failed process
	message string
	// ResponseCode represents the HTTP status code the error will return
	responseCode int
}

// Error returns the error in a string format
func (a apierror) Error() string {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(a)
	return buf.String()
}

// StatusCode set the http status code of the response
func (a apierror) StatusCode() int {
	return a.responseCode
}

// MarshalJSON defines the json representation of the error
func (a apierror) MarshalJSON() ([]byte, error) {
	vals := map[string]interface{}{}
	vals["error"] = map[string]interface{}{
		"message": a.message,
		"code":    a.code,
	}
	return json.Marshal(vals)
}

// NotFound returns a not found error
func NotFound(code string, err error) error {
	return apierror{
		code:         code,
		responseCode: http.StatusNotFound,
		message:      err.Error(),
	}
}

// InvalidRequest returns a bad request error
func InvalidRequest(code string, err error) error {
	return apierror{
		code:         code,
		responseCode: http.StatusBadRequest,
		message:      err.Error(),
	}
}

// Internal returns an internal server error
func Internal(code string, err error) error {
	return apierror{
		code:         code,
		responseCode: http.StatusInternalServerError,
		message:      err.Error(),
	}
}

// Unauthorized returns an unauthorized error
func Unauthorized(code string, err error) error {
	return apierror{
		code:         code,
		message:      err.Error(),
		responseCode: http.StatusUnauthorized,
	}
}
