package raw

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// func parseHeadEtag(rep *http.Response) (string, error) {
// 	etag := rep.Header.Get("Etag")
// 	if etag == "" {
// 		return "", fmt.Errorf(
// 			"'Etag' missing from response headers",
// 		)
// 	}
// 	return etag, nil
// }

func parseHeadExpires(rep *http.Response) (time.Time, error) {
	datestring := rep.Header.Get("Expires")
	if datestring == "" {
		return time.Time{}, fmt.Errorf(
			"'Expires' missing from response headers",
		)
	}
	date, err := time.Parse(time.RFC1123, datestring)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"error parsing 'Expires' header: %w",
			err,
		)
	}
	return date, nil
}

func parseHeadPages(rep *http.Response) (int, error) {
	pagesstring := rep.Header.Get("X-Pages")
	if pagesstring == "" {
		return 0, fmt.Errorf(
			"'X-Pages' missing from response headers",
		)
	}
	pages, err := strconv.Atoi(pagesstring)
	if err != nil {
		return 0, fmt.Errorf(
			"error parsing 'X-Pages' header: %w",
			err,
		)
	}
	return pages, nil
}
