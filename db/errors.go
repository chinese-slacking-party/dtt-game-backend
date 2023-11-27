package db

import (
	driver "go.mongodb.org/mongo-driver/mongo"
)

type ErrDuplicateKey struct {
	Internal error
}

func (e *ErrDuplicateKey) Error() string {
	if e == nil || e.Internal == nil {
		return "ErrDuplicateKey(nil)"
	}
	return string(e.Internal.Error())
}

var (
	ErrNotFound = driver.ErrNoDocuments
)
