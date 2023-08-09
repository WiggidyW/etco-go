package rawclient

import (
	"fmt"
	"io"
	"net/http"
)

type CacheCreateError struct{ error }
type CacheGetError struct{ error }
type CacheSetError struct{ error }
type RequestParamsError struct{ error }
type HttpError struct{ error }
type MalformedResponseBody struct{ error }
type MalformedResponseHeaders struct{ error }

type StatusError struct {
	Url      string
	Code     int
	CodeText string
	EsiText  string
}

func (e StatusError) Error() string {
	errstr := fmt.Sprintf(
		"ESI Server Request '%s' returned Response Code '%s'",
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

func newStatusError(rep *http.Response) StatusError {
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
