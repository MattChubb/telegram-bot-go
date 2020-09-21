package markov

import (
	"github.com/mb-14/gomarkov"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Brain struct {
    chain *gomarkov.Chain
}

func (brain *Brain) Init(order int) {
	brain.chain = gomarkov.NewChain(order)
    log.Debug("Braindump: ", brain)
}

func (brain *Brain) Train(data string) error {
    log.Debug("Braindump: ", brain)
    log.Debug("Training data: ", data)
    processedData := processString(data)
    log.Debug("Processed into: ", processedData)
    brain.chain.Add(processedData)
    return nil
}

func Generate() {}

func processString(rawString string) []string {
	//TODO Handle punctuation other than spaces
	return strings.Split(strings.ToLower(rawString), " ")
}
