package doublemarkov

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
    "github.com/TwinProduction/go-away"
	log "github.com/sirupsen/logrus"
    "regexp"
	"strings"
    common "github.com/MattChubb/telegram-bot-go/brain"
)

//TODO Use a bi-directional markov chain instead of 2 separate chains to lower memory footprint
type Brain struct {
    bckChain    *gomarkov.Chain
    fwdChain    *gomarkov.Chain
    lengthLimit int
}

type brainJSON struct {
    BckChain    *gomarkov.Chain
    FwdChain    *gomarkov.Chain
    LengthLimit int
}

func (brain Brain) MarshalJSON() ([]byte, error) {
	log.Info("Saving chain...")

    obj := brainJSON{
        brain.bckChain,
        brain.fwdChain,
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

    brain.bckChain = obj.BckChain
    brain.fwdChain = obj.FwdChain
    brain.lengthLimit = obj.LengthLimit
    log.Debug("Braindump: ", brain)

    return nil
}

func (brain *Brain) Init(order int, lengthLimit int) {
	brain.bckChain = gomarkov.NewChain(order)
	brain.fwdChain = gomarkov.NewChain(order)
    brain.lengthLimit = lengthLimit
    log.Debug("Braindump: ", brain)
}

func (brain *Brain) Train(data string) error {
    log.Debug("Braindump: ", brain)
    log.Debug("Training data: ", data)

    processedData := common.ProcessString(data)
    log.Debug("Processed into: ", processedData)

    brain.fwdChain.Add(processedData)
    reverse(processedData)
    log.Debug("Reversed: ", processedData)
    brain.bckChain.Add(processedData)
    return nil
}

func (brain *Brain) Generate(prompt string) (string, error) {
    processedPrompt := common.ProcessString(prompt)
	subject := []string{}
	if len(processedPrompt) > 0 {
		subject = common.ExtractSubject(processedPrompt)
	}
	//TODO Any other clever Markov hacks?
	sentence := brain.generateSentence(brain.bckChain, subject)
	end := brain.generateSentence(brain.fwdChain, subject)

    sentence = sentence[1:]
    reverse(sentence)
    sentence = append(sentence, end...)
    sentence[0] = strings.Title(sentence[0])

    return strings.Join(sentence, " "), nil
}

func (brain *Brain) generateSentence(chain *gomarkov.Chain, init []string) []string {
    // The length of our initialisation chain needs to match the Markov order
	tokens := []string{}
    order := chain.Order
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
		len(tokens) < brain.lengthLimit { //TODO lengthlimit should apply to the whole generated sentence instead of the individual halves
		next, err := chain.Generate(tokens[(len(tokens) - 1):])
		if err != nil {
            if match, err := regexp.Match(`Unknown ngram.*`, []byte(err.Error())); !match {
			    log.Fatal(err)
            }
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

func reverse(ss []string) {
    last := len(ss) - 1
    for i := 0; i < len(ss)/2; i++ {
        ss[i], ss[last-i] = ss[last-i], ss[i]
    }
}
