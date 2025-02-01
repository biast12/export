package validator

import (
	"archive/zip"
	"crypto/ed25519"
	"github.com/TicketsBot/data-self-service/internal/utils"
	"io"
)

func (v *Validator) validateSignature(zipReader *zip.Reader, fileName string, data []byte) (int64, error) {
	f, err := zipReader.Open(fileName + ".sig")
	if err != nil {
		return 0, err
	}

	signature, err := io.ReadAll(v.newLimitReader(f))
	if err != nil {
		return 0, err
	}

	decoded, err := utils.Base64Decode(string(signature))
	if err != nil {
		return 0, err
	}

	if !ed25519.Verify(v.publicKey, data, decoded) {
		return 0, ErrValidationFailed
	}

	return int64(len(signature)), nil
}
