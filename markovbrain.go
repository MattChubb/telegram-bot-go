package markovbrain

import (
	"github.com/mb-14/gomarkov"
)

type Brain struct {
    chain gomarkov.Chain
}

func (brain Brain) Init(order int) {
	brain.chain = *gomarkov.NewChain(order)
}

func Train() {}

func Generate() {}
