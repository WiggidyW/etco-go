package esi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/WiggidyW/etco-go/error/esierror"
)

func parseHeadExpires(rep *http.Response) (
	expires time.Time,
	err error,
) {
	datestring := rep.Header.Get("Expires")
	if datestring == "" {
		err = esierror.MalformedResponseHeaders{
			Err: fmt.Errorf("'Expires' missing from response headers"),
		}
		return expires, err
	}

	expires, err = time.Parse(time.RFC1123, datestring)
	if err != nil {
		err = esierror.MalformedResponseHeaders{
			Err: fmt.Errorf(
				"error parsing 'Expires' header: %w",
				err,
			),
		}
	}
	return expires, err
}

func parseHeadPages(rep *http.Response) (
	pages int,
	err error,
) {
	pagesstring := rep.Header.Get("X-Pages")
	if pagesstring == "" {
		err = esierror.MalformedResponseHeaders{
			Err: fmt.Errorf("'X-Pages' missing from response headers"),
		}
		return pages, err
	}

	pages, err = strconv.Atoi(pagesstring)
	if err != nil {
		err = esierror.MalformedResponseHeaders{
			Err: fmt.Errorf(
				"error parsing 'X-Pages' header: %w",
				err,
			),
		}
	}
	return pages, err
}
