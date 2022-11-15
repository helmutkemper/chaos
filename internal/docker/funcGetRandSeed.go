package docker

import (
	"math/rand"
	"time"
)

// getRandSeed
//
// English:
//
//	Generate random number seed
//
//	 Output:
//	   seed: random number seed
//
// Português:
//
//	Gera a semente do número aleatório
//
//	 Saída:
//	   seed: semente do número aleatório
func (e *ContainerBuilder) getRandSeed() (seed *rand.Rand) {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}
