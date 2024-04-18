package uuid

import (
	guuid "github.com/google/uuid"
)

type NewUuidFunc func() string

var NewUuid NewUuidFunc

func ResetUuidImplementation() {
	NewUuid = func() string {
		return guuid.NewString()
	}
}

func Uuid() string {
	return guuid.NewString()
}
