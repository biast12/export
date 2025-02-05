package validator

import (
	"archive/zip"
	"encoding/json"
	"github.com/TicketsBot/data-self-service/internal/model/dto"
	"io"
)

func (v *Validator) ValidateGuildData(input io.ReaderAt, size int64) (*dto.GuildData, error) {
	reader, err := zip.NewReader(input, size)
	if err != nil {
		return nil, err
	}

	f, err := reader.Open("data.json")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(v.newLimitReader(f))
	if err != nil {
		return nil, err
	}

	if _, err := v.validateSignature(reader, "data.json", data); err != nil {
		return nil, err
	}

	var guildData dto.GuildData
	if err := json.Unmarshal(data, &guildData); err != nil {
		return nil, err
	}

	return &guildData, nil
}
