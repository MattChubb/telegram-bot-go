package markov

import (
	"github.com/mb-14/gomarkov"
    "github.com/TwinProduction/go-away"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

type Brain struct {
    chain       *gomarkov.Chain
    lengthLimit int
}

func (brain *Brain) Init(order int, lengthLimit int) {
	brain.chain = gomarkov.NewChain(order)
    brain.lengthLimit = lengthLimit
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

func (brain *Brain) Generate(prompt string) (string, error) {
    processedPrompt := processString(prompt)
	subject := []string{}
	if len(processedPrompt) > 0 {
		subject = extractSubject(processedPrompt)
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

func processString(rawString string) []string {
	//TODO Handle punctuation other than spaces
	return strings.Split(strings.ToLower(rawString), " ")
}

func extractSubject(message []string) []string {
    //TODO Do something cleverer with subject extraction
    //TODO Extract more than one word as a subject
    trimmedMessage := trimMessage(message)
    subject := []string{}
    if len(trimmedMessage) > 0 {
        subject = append(subject, trimmedMessage[rand.Intn(len(trimmedMessage))])
    }
    return subject
}

func trimMessage(message []string) []string {
    trimmedMessage := []string{}
    for _, word := range message {
        //TODO Only exclude self-mentions
        if ! isStopWord(word) && word[0] != '@' {
            trimmedMessage = append(trimmedMessage, word)
        }
    }
    return trimmedMessage
}

func isStopWord(word string) bool {
    stopWords := []string{"the", "and", "to", "a", "i", "in", "be", "of", "that", "have", "it", }
    for _, stopWord := range stopWords {
        if word == stopWord {
            return true
        }
    }

    return false
}
