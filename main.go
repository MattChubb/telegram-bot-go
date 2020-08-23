package main

import (
//    tb "gopkg.in/tucnak/telebot.v2"
    "github.com/mb-14/gomarkov"
    "fmt"
    "strings"
)

func main() {
    //Initilise Telegram bot

    //Initialise chain
    chain := gomarkov.NewChain(2)

    //Open source data dir

    //Train markov chain
    chain.Add(processString("test data test data test data"))
    chain.Add(processString("test data test data test data"))
    chain.Add(processString("test data test data test data"))
    next, _ := chain.Generate([]string{"test", "data", "data"})
    fmt.Println(next)

    //Connect to telegram
}

func processString(rawString string) []string {
    //TODO Handle punctuation other than spaces
    return strings.Split(rawString, " ")
}
