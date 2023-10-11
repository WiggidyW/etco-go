package builderenv

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateFileAndValidateCreateFileDirectory(filePath string) (err error) {
	dirPath := filepath.Dir(filePath)
	if err := validateCreateDirectory(dirPath); err != nil {
		return err
	}
	if err := validateFile(filePath); err != nil {
		return err
	}
	return nil
}

// validates the directory, creating it if needed and returning error if
// unable to create it or the path is invalid
func validateCreateDirectory(dirPath string) (err error) {
	dirInfo, err := os.Stat(dirPath)

	if err != nil {

		if os.IsNotExist(err) {
			// create the missing dirPath
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				return err
			}

		} else {
			// error getting dir info
			return err
		}

	} else if !dirInfo.IsDir() {
		// return an error if the dirPath exists and is not a directory
		return fmt.Errorf(
			"dir path %s is a file and not a directory",
			dirPath,
		)
	}

	return nil
}

// simply returns an error if the file path is a directory
func validateFile(filePath string) (err error) {
	fileInfo, err := os.Stat(filePath)

	if err != nil {

		if os.IsNotExist(err) {
			// acceptable error, file does not exist
			return nil

		} else {
			// error getting file info
			return err
		}

	} else if fileInfo.IsDir() {
		// return an error if the filePath is a directory
		return fmt.Errorf(
			"file path %s is a directory and not a file",
			filePath,
		)
	}

	return nil
}
