package doublemarkov

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
	log "github.com/sirupsen/logrus"
    "math"
    "regexp"
	"strings"
    common "github.com/MattChubb/telegram-bot-go/brain"
    markov "github.com/MattChubb/telegram-bot-go/brain/markov"
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
		subject = common.ExtractSubject(processedPrompt, brain.fwdChain.Order)
	}
	//TODO Any other clever Markov hacks?
	sentence := brain.generateSentence(brain.bckChain, subject)
	end := brain.generateSentence(brain.fwdChain, subject)

    if len(sentence) > brain.bckChain.Order {
        // Don't start a sentence with punctuation
        if match, _ := regexp.Match(`\W`, []byte(sentence[len(sentence)-1])); match {
            sentence = sentence[brain.bckChain.Order:len(sentence)-1]
        } else {
            sentence = sentence[brain.bckChain.Order:]
        }
    } else {
        //Sentence is just the subject, which is duplicated in the fwd chain
        sentence = []string{}
    }
    reverse(sentence)
    sentence = append(sentence, end...)
    sentence[0] = strings.Title(sentence[0])

    return strings.Join(sentence, ""), nil
}

func (brain *Brain) generateSentence(chain *gomarkov.Chain, init []string) []string {
    log.Debug("Input: ", init)
    order := chain.Order
    tokens := markov.GenerateInitialToken(init, order)
    log.Debug("Initial token: ", tokens)

	for tokens[len(tokens)-1] != gomarkov.EndToken &&
        len(tokens) < int(math.Round(float64(brain.lengthLimit)/2)) {
        next := markov.GenerateNextToken(chain, tokens)
        tokens = append(tokens, next)
	}

	//Don't include the start or end token in our response
    return markov.TrimTokens(tokens)
}

func reverse(ss []string) {
    last := len(ss) - 1
    for i := 0; i < len(ss)/2; i++ {
        ss[i], ss[last-i] = ss[last-i], ss[i]
    }
}
