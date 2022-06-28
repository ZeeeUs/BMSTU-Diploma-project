package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/google/uuid"
)

type File struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Bytes []byte `json:"file"`
}

// CreateFileUnit is data transfer object
type CreateFileUnit struct {
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}

func (u CreateFileUnit) NormalizeName() {
	u.PayloadName = strings.ReplaceAll(u.PayloadName, " ", "_")
}

func NewFile(dto CreateFileUnit) (*File, error) {
	bytes, err := ioutil.ReadAll(dto.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create file model. err: %w", err)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate file id. err: %w", err)
	}

	return &File{
		ID:    id.String(),
		Name:  dto.PayloadName,
		Size:  dto.PayloadSize,
		Bytes: bytes,
	}, nil
}
