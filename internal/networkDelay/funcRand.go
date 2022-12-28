package networkdelay

import (
	"math/rand"
	"time"
)

func (e *Proxy) rand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
