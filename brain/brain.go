package brain

import (
    "regexp"
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
	return regexp.MustCompile(`\b`).Split(strings.ToLower(rawString), -1)
}

func ExtractSubject(message []string, length int) []string {
    trimmedMessage := trimMessage(message)
    var subject string
    if len(trimmedMessage) > 0 {
        subject = trimmedMessage[rand.Intn(len(trimmedMessage))]
    } else {
        //If there's nothing but stopwords, return nothing
        return []string{}
    }

    if length == 1 {
        //Short-circuit as we don't need to pad
        return []string{subject}
    }

    subjectWords := []string{}
    for _, word := range message {
        subjectWords = append(subjectWords, word)
        if len(subjectWords) > length {
            //We want the main subject word to be roughly halfway through the
            //subject words, or at the beginning if that's not possible

            if subjectWords[0] == subject {
                subjectWords = subjectWords[:len(subjectWords)-1]
                break
            }

            subjectWords = subjectWords[1:]
            if subjectWords[length/2] == subject {
                break
            }
        }
    }

    return subjectWords
}

func trimMessage(message []string) []string {
    trimmedMessage := []string{}
    for _, word := range message {
        //TODO Only exclude self-mentions
        if  match, _ := regexp.Match(`\W`, []byte(word)); ! match && len(word) > 0 && ! isStopWord(word) && word[0] != '@' {
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
