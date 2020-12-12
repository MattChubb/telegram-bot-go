package brain

import (
	"strings"
	"math/rand"
)

type Brain interface{
    //TODO Add a more flexible init method
    Init(o int, l int)
    Train(d string) error
    Generate(p string) (string, error)
}

func ProcessString(rawString string) []string {
	//TODO Handle punctuation other than spaces
	return strings.Split(strings.ToLower(rawString), " ")
}

func ExtractSubject(message []string) []string {
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
        if len(word) > 0 && ! isStopWord(word) && word[0] != '@' {
            trimmedMessage = append(trimmedMessage, word)
        }
    }
    return trimmedMessage
}

func isStopWord(word string) bool {
    for _, stopWord := range stopWords {
        if word == stopWord {
            return true
        }
    }

    return false
}
