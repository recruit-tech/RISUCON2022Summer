package random

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func init() {
	uuid.EnableRandPool()
	rand.Seed(time.Now().UnixNano())
}

func ID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	return uuid.String()
}

func Duration(min, max time.Duration) time.Duration {
	return time.Duration(rand.Intn(int(max-min))) + min
}
