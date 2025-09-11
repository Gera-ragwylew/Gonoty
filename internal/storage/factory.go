package storage

import (
	"fmt"
)

type Type string

const (
	Redis       Type = "redis"
	Postgres    Type = "postgres"
	Mock        Type = "mock"
	FileStorage Type = "file storage"
)

func NewStorage(storageType Type) (Storage, error) {
	switch storageType {
	// case Redis:
	// 	return redisstorage.NewRedisStorage(), nil
	case Postgres:
		return nil, nil
	// case Mock:
	// 	return mock.NewMockStorage(), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}
