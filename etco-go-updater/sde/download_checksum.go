package sde

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/WiggidyW/chanresult"
)

const (
	CHECKSUM_URL            = "https://eve-static-data-export.s3-eu-west-1.amazonaws.com/tranquility/checksum"
	CHECKSUM_CAPACITY_BYTES = 32
	CHECKSUM_TIMEOUT        = 60 * time.Second
)

var RE_VALID_CHECKSUM = regexp.MustCompile(`^[0-9a-f]{32}$`)

func TransceiveDownloadChecksum(
	ctx context.Context,
	httpClient *http.Client,
	userAgent string,
	chnSend chanresult.ChanSendResult[string],
) error {
	checksum, err := DownloadChecksum(ctx, httpClient, userAgent)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(checksum)
	}
}

func DownloadChecksum(
	ctx context.Context,
	httpClient *http.Client,
	userAgent string,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, CHECKSUM_TIMEOUT)
	defer cancel()

	// initialize a buffer to store the checksum file
	buf := bytes.NewBuffer(make([]byte, 0, CHECKSUM_CAPACITY_BYTES))

	// prepare the request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		CHECKSUM_URL,
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	// download the checksum file
	rep, err := httpClient.Do(req)
	if err != nil {
		return "", err
	} else if rep.StatusCode != http.StatusOK {
		rep.Body.Close()
		return "", fmt.Errorf("bad status code: %d", rep.StatusCode)
	}

	// copy the checksum file into the buffer
	_, err = io.Copy(buf, rep.Body)
	rep.Body.Close()
	if err != nil {
		return "", err
	}

	// verify that the checksum is valid
	checksum := buf.String()
	if !RE_VALID_CHECKSUM.MatchString(checksum) {
		return "", fmt.Errorf("invalid checksum: %s", checksum)
	}

	return checksum, nil
}
