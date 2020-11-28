package markov

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
    "github.com/TwinProduction/go-away"
	log "github.com/sirupsen/logrus"
	"strings"
    common "github.com/MattChubb/telegram-bot-go/brain"
)

type Brain struct {
    chain       *gomarkov.Chain
    lengthLimit int
}

type brainJSON struct {
    Chain       *gomarkov.Chain
    LengthLimit int
}

func (brain Brain) MarshalJSON() ([]byte, error) {
	log.Info("Saving chain...")

    obj := brainJSON{
        brain.chain,
        brain.lengthLimit,
    }

    return json.Marshal(obj)
}

func (brain *Brain) UnmarshalJSON(b []byte) error {
	var obj brainJSON
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}

    brain.lengthLimit = obj.LengthLimit
    brain.chain = obj.Chain
    log.Debug("Braindump: ", brain)

    return nil
}

func (brain *Brain) Init(order int, lengthLimit int) {
	brain.chain = gomarkov.NewChain(order)
    brain.lengthLimit = lengthLimit
    log.Debug("Braindump: ", brain)
}

func (brain *Brain) Train(data string) error {
    log.Debug("Braindump: ", brain)
    log.Debug("Training data: ", data)
    processedData := common.ProcessString(data)
    log.Debug("Processed into: ", processedData)
    brain.chain.Add(processedData)
    return nil
}

func (brain *Brain) Generate(prompt string) (string, error) {
    processedPrompt := common.ProcessString(prompt)
	subject := []string{}
	if len(processedPrompt) > 0 {
		subject = common.ExtractSubject(processedPrompt)
	}
	//TODO Bi-directional generation using both a forwards and a backwards trained Markov chains
	//TODO Any other clever Markov hacks?
	sentence := brain.generateSentence(subject)
    sentence[0] = strings.Title(sentence[0])
    return strings.Join(sentence, " "), nil
}

func (brain *Brain) generateSentence(init []string) []string {
	// This function has been separated from response generation to allow bidirectional generation later

    //Train on the initial tokens to avoid unknown n-grams
    brain.chain.Add(init)

    // The length of our initialisation chain needs to match the Markov order
	tokens := []string{}
    order := brain.chain.Order
	if len(init) < order {
		for i := 0; i < order; i++ {
			tokens = append(tokens, gomarkov.StartToken)
		}
		tokens = append(tokens, init...)
	} else if len(init) > order {
		tokens = init[:order]
	} else {
		tokens = init
	}

	for tokens[len(tokens)-1] != gomarkov.EndToken &&
		len(tokens) < brain.lengthLimit {
		next, err := brain.chain.Generate(tokens[(len(tokens) - 1):])
		if err != nil {
			log.Fatal(err)
		}

        //TODO Implement a replacement wordfilter instead of just removing profanity
		if len(next) > 0 && ! goaway.IsProfane(next) {
			tokens = append(tokens, next)
		} else {
			tokens = append(tokens, gomarkov.EndToken)
		}
	}

	//Don't include the start or end token in our response
	tokens = tokens[:len(tokens)-1]
	if tokens[0] == gomarkov.StartToken {
		tokens = tokens[1:]
	}
	return tokens
}
