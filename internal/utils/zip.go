package utils

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"github.com/TicketsBot/export/internal/config"
	"io"
)

func BuildZip(cfg config.WorkerConfig, files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer

	w := zip.NewWriter(&buf)
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, cfg.Daemon.CompressionLevel)
	})

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
