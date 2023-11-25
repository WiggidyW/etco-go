package esierror

import (
	"fmt"
	"io"
	"net/http"
)

type RequestParamsError struct{ Err error }

func (e RequestParamsError) Unwrap() error { return e.Err }
func (e RequestParamsError) Error() string {
	return fmt.Sprintf("RequestParamsError: %s", e.Err)
}

type HttpError struct{ Err error }

func (e HttpError) Unwrap() error { return e.Err }
func (e HttpError) Error() string {
	return fmt.Sprintf("HttpError: %s", e.Err)
}

type MalformedResponseBody struct{ Err error }

func (e MalformedResponseBody) Unwrap() error { return e.Err }
func (e MalformedResponseBody) Error() string {
	return fmt.Sprintf("MalformedResponseBody: %s", e.Err)
}

type MalformedResponseHeaders struct{ Err error }

func (e MalformedResponseHeaders) Unwrap() error { return e.Err }
func (e MalformedResponseHeaders) Error() string {
	return fmt.Sprintf("MalformedResponseHeaders: %s", e.Err)
}

type AuthRefreshError struct{ Err error }

func (e AuthRefreshError) Unwrap() error { return e.Err }
func (e AuthRefreshError) Error() string {
	return fmt.Sprintf("AuthRefreshError: %s", e.Err)
}

// req refresh token != rep refresh token
// This means that CCP has enabled rotation for web refresh tokens
// I think they'll never do this, but if they do, we'll need to handle it
type AuthRefreshMismatch struct{ App string }

func (e AuthRefreshMismatch) Unwrap() error { return nil }
func (e AuthRefreshMismatch) Error() string {
	return fmt.Sprintf("(CRITICAL, VERY BAD) AuthRefreshMismatch: %s", e.App)
}

type StatusError struct {
	Url      string
	Code     int
	CodeText string
	EsiText  string
}

func NewStatusError(rep *http.Response) StatusError {
	var body_str string
	body_bytes, err := io.ReadAll(rep.Body)
	if err != nil {
		body_str = ""
	} else {
		body_str = string(body_bytes)
	}
	return StatusError{
		Url:      rep.Request.URL.String(),
		Code:     rep.StatusCode,
		CodeText: rep.Status,
		EsiText:  body_str,
	}
}

func (e StatusError) Unwrap() error { return nil }
func (e StatusError) Error() string {
	errstr := fmt.Sprintf(
		"StatusError: '%s' returned Response Code '%s'",
		e.Url,
		e.CodeText,
	)
	if e.EsiText == "" {
		errstr += " with no error message"
	} else {
		errstr += fmt.Sprintf(
			" with error message '%s'",
			e.EsiText,
		)
	}
	return errstr
}
