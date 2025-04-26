package utils

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func GenUUID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
