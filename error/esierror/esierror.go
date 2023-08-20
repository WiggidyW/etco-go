package esierror

import (
	"fmt"
	"io"
	"net/http"
)

type RequestParamsError struct{ Err error }

func (e RequestParamsError) Error() string {
	return fmt.Sprintf("RequestParamsError: %s", e.Err)
}

type HttpError struct{ Err error }

func (e HttpError) Error() string {
	return fmt.Sprintf("HttpError: %s", e.Err)
}

type MalformedResponseBody struct{ Err error }

func (e MalformedResponseBody) Error() string {
	return fmt.Sprintf("MalformedResponseBody: %s", e.Err)
}

type MalformedResponseHeaders struct{ Err error }

func (e MalformedResponseHeaders) Error() string {
	return fmt.Sprintf("MalformedResponseHeaders: %s", e.Err)
}

type StatusError struct {
	Url      string
	Code     int
	CodeText string
	EsiText  string
}

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
