package markov

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
    "github.com/TwinProduction/go-away"
	log "github.com/sirupsen/logrus"
    "regexp"
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
    log.Debug("Input: ", prompt)
    processedPrompt := common.ProcessString(prompt)
    log.Debug("Processed into: ", processedPrompt)

	subject := []string{}
	if len(processedPrompt) > 0 {
		subject = common.ExtractSubject(processedPrompt, brain.chain.Order)
	}
	//TODO Any other clever Markov hacks?
	sentence := brain.generateSentence(subject)
    sentence[0] = strings.Title(sentence[0])
    return strings.Join(sentence, ""), nil
}

func (brain *Brain) generateSentence(init []string) []string {
    log.Debug("Input: ", init)
    order := brain.chain.Order
    tokens := GenerateInitialToken(init, order)
    log.Debug("Initial token: ", tokens)

	for tokens[len(tokens)-1] != gomarkov.EndToken &&
		len(tokens) < brain.lengthLimit {
        next := GenerateNextToken(brain.chain, tokens)
        tokens = append(tokens, next)
	}

	//Don't include the start or end token in our response
    return TrimTokens(tokens)
}

//Exported so that DoubleMarkov can also use it
//It could also go in brain, but it's markov-specific
func GenerateInitialToken(init []string, order int) []string {
	tokens := []string{}
    // The length of our initialisation chain needs to match the Markov order
	if len(init) < order {
		for i := 0; i < order - len(init) ; i++ {
			tokens = append(tokens, gomarkov.StartToken)
		}
		tokens = append(tokens, init...)
	} else if len(init) > order {
		tokens = init[:order]
	} else {
		tokens = init
	}

    return tokens
}

func TrimTokens(tokens []string) []string {
	tokens = tokens[:len(tokens)-1]
	for tokens[0] == gomarkov.StartToken {
		tokens = tokens[1:]
	}
	return tokens
}

func GenerateNextToken(chain *gomarkov.Chain, tokens []string) string {
    next, err := chain.Generate(tokens[(len(tokens) - chain.Order):])
    if err != nil {
        errormsg := err.Error()
        if match, _ := regexp.Match(`Unknown ngram.*`, []byte(errormsg)); !match {
            log.Fatal("Error generating from Markov chain: ", errormsg)
        }
    }

    //TODO Implement a replacement wordfilter instead of just removing profanity
    if len(next) > 0 && ! goaway.IsProfane(next) {
        return next
    } else {
        return gomarkov.EndToken
    }
}
