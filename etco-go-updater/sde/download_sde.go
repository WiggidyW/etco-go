package sde

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SDE_URL     = "https://eve-static-data-export.s3-eu-west-1.amazonaws.com/tranquility/sde.zip"
	SDE_TIMEOUT = 1800 * time.Second
)

func downloadSDE(
	ctx context.Context,
	httpClient *http.Client,
	userAgent string,
	pathSDE string,
) error {
	ctx, cancel := context.WithTimeout(ctx, SDE_TIMEOUT)
	defer cancel()

	// verify that dst is a directory, or create it
	if fi, err := os.Stat(pathSDE); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(pathSDE, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf(
			"destination is not a directory: %s",
			pathSDE,
		)
	}

	// initialize a temporary zip file to store the sde
	zipFile, err := os.CreateTemp("", "sde-*.zip")
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
		_ = os.Remove(zipFile.Name())
	}()

	// prepare the request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		SDE_URL,
		nil,
	)
	if err != nil {
		return err
	}

	// download the sde
	rep, err := httpClient.Do(req)
	if err != nil {
		return err
	} else if rep.StatusCode != http.StatusOK {
		rep.Body.Close()
		return fmt.Errorf("bad status code: %d", rep.StatusCode)
	}

	// copy the sde into the temporary zip file
	_, err = io.Copy(zipFile, rep.Body)
	rep.Body.Close()
	if err != nil {
		return err
	}

	// extract the contents into dstPath
	if err := extractSDE(zipFile, pathSDE); err != nil {
		return err
	} else {
		return nil
	}
}

func extractSDE(zipFile *os.File, targetDir string) error {
	// Seek to the beginning of the file before reading
	_, err := zipFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// create a zip reader for the zip file
	var r *zip.Reader
	if stat, err := zipFile.Stat(); err != nil {
		return err
	} else {
		r, err = zip.NewReader(zipFile, stat.Size())
		if err != nil {
			return err
		}
	}

	// Ensure that there's only one of three valid root level directories
	for _, f := range r.File {
		if !strings.HasPrefix(f.Name, "bsd/") &&
			!strings.HasPrefix(f.Name, "fsd/") &&
			!strings.HasPrefix(f.Name, "universe/") {
			return errors.New(
				"invalid entry in the zip: " + f.Name,
			)
		}
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create a directory if needed
		fpath := filepath.Join(targetDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Ensure the directory for the file exists
			if err = os.MkdirAll(
				filepath.Dir(fpath),
				os.ModePerm,
			); err != nil {
				return err
			}

			// Create the file
			outFile, err := os.OpenFile(
				fpath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				f.Mode(),
			)
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
