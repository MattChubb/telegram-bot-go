package main

import (
//    tb "gopkg.in/tucnak/telebot.v2"
    "github.com/mb-14/gomarkov"
    "fmt"
    "strings"
)

const tokensLengthLimit = 32

func main() {
    //Initilise Telegram bot

    //Initialise chain
    chain := gomarkov.NewChain(2)

    //Open source data dir

    //Train markov chain
    chain.Add(processString("test data test data test data"))
    chain.Add(processString("data test data test data"))
    chain.Add(processString("test data test data test data"))

    sentence := generateSentence(chain, []string{"test"})
    fmt.Println(sentence)

    //Connect to telegram
}

func processString(rawString string) []string {
    //TODO Handle punctuation other than spaces
    return strings.Split(rawString, " ")
}

func generateSentence(chain *gomarkov.Chain, init []string) []string {
    tokens := []string{gomarkov.StartToken}
    tokens = append(tokens, init...)

    for tokens[len(tokens) - 1] != gomarkov.EndToken &&
        len(tokens) < tokensLengthLimit {
        next, _ := chain.Generate(tokens[(len(tokens) - 1):] )
        tokens = append(tokens, next)
    }

    return tokens
}
