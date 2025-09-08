package storage

import (
	mock "Gonoty/internal/storage/test_mock"
	"fmt"
)

type Type string

const (
	Postgres    Type = "postgres"
	Mock        Type = "mock"
	FileStorage Type = "file storage"
)

func NewStorage(storageType Type) (Storage, error) {
	switch storageType {
	case Postgres:
		return nil, nil
	case Mock:
		return mock.NewMockStorage(), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
