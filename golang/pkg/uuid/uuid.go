package uuid

import (
	guuid "github.com/google/uuid"
)

type NewUuidFunc func() string

var NewUuid NewUuidFunc = func() string {
	return guuid.NewString()
}

func ResetUuidImplementation() {
	NewUuid = func() string {
		return guuid.NewString()
	}
}

func Uuid() string {
	return NewUuid()
}
