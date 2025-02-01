package utils

import (
	"archive/zip"
	"bytes"
)

func BuildZip(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer

	w := zip.NewWriter(&buf)

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return nil, err
		}

		if _, err := f.Write(content); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
