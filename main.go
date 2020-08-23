package main

import (
//    tb "gopkg.in/tucnak/telebot.v2"
    "github.com/mb-14/gomarkov"
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
)

const (
    tokensLengthLimit = 32
    order = 1
)

func main() {
    //Initialise
    //Initilise Telegram bot

    //Initialise chain
    chain := gomarkov.NewChain(order)

    //Open source data dir
    source_file, err := os.Open("./source_data/data")
    if err != nil {
        log.Fatal(err)
    }
    defer source_file.Close()
    scanner := bufio.NewScanner(source_file)


    //Train
    //Train markov chain
    for scanner.Scan() {
        chain.Add(processString(scanner.Text()))
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }


    //Connect Markov to Telegram
    //Process input

    //Identify whether to respond

    //Generate response
    sentence := generateSentence(chain, []string{"I", "think"})

    //Respond with generated response
    fmt.Println(sentence)

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
       // fmt.Println(next)
        if len(next) > 0 {
            tokens = append(tokens, next)
        } else {
            tokens = append(tokens, gomarkov.EndToken)
        }
    }

    return tokens
}
