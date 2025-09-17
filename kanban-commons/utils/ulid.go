package utils

import (
	"crypto/rand"
	"math"
	"time"

	"github.com/oklog/ulid/v2"
)

type Ulid struct {
	entropy *ulid.MonotonicEntropy
}

// Call once before using CreateID
func NewUlidGenerator() *Ulid {
	seed := ulid.Monotonic(rand.Reader, math.MaxInt64)
	return &Ulid{entropy: seed}
}

func (u *Ulid) CreateID() string {
	t := time.Now().UTC()
	id := ulid.MustNew(ulid.Timestamp(t), u.entropy)
	return id.String()
}
