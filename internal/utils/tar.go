package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"time"
)

func BuildTarball(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer

	writer := tar.NewWriter(&buf)

	dirNames := make(map[string]struct{})
	for name, content := range files {
		parts := bytes.Split([]byte(name), []byte("/"))

		for i := 1; i < len(parts); i++ {
			dir := parts[:i]
			dirName := string(bytes.Join(dir, []byte("/"))) + "/"

			if _, ok := dirNames[dirName]; ok {
				continue
			} else {
				dirNames[dirName] = struct{}{}
			}

			header := &tar.Header{
				Name:     dirName,
				Typeflag: tar.TypeDir,
				Mode:     0644,
				ModTime:  time.Now(),
			}

			if err := writer.WriteHeader(header); err != nil {
				return nil, err
			}
		}

		header := &tar.Header{
			Name:    name,
			Size:    int64(len(content)),
			Mode:    0644,
			ModTime: time.Now(),
		}

		if err := writer.WriteHeader(header); err != nil {
			return nil, err
		}

		if _, err := writer.Write(content); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func MergeTarballs(tarballs ...io.Reader) ([]byte, error) {
	var buf bytes.Buffer

	writer := tar.NewWriter(&buf)

	for _, tarball := range tarballs {
		reader := tar.NewReader(tarball)
		for {
			header, err := reader.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return nil, err
			}

			if err := writer.WriteHeader(header); err != nil {
				return nil, err
			}

			if _, err := io.Copy(writer, reader); err != nil {
				return nil, err
			}
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TarballToZip(tarball io.Reader) ([]byte, error) {
	var buf bytes.Buffer

	writer := zip.NewWriter(&buf)

	reader := tar.NewReader(tarball)
	for {
		header, err := reader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		if header.Typeflag == tar.TypeReg {
			file, err := writer.Create(header.Name)
			if err != nil {
				return nil, err
			}

			if _, err := io.Copy(file, reader); err != nil {
				return nil, err
			}
		} else {
			continue
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
