package loader_

import (
	"encoding/gob"
	"os"
)

func gobFsLoad[T any](t *T, fsPath string) error {
	// open file
	f, err := os.Open(fsPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// decode
	dec := gob.NewDecoder(f)
	err = dec.Decode(t)
	if err != nil {
		return err
	}

	return nil
}
