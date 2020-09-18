package markov

import (
	"github.com/mb-14/gomarkov"
)

type Brain struct {
    chain gomarkov.Chain
}

func (brain Brain) Init(order int) {
	brain.chain = *gomarkov.NewChain(order)
}

func (brain Brain) Train(data string) error {
    return nil
}

func Generate() {}
