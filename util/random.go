package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInit(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) 
}
