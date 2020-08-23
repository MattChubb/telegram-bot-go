package main

import (
    "github.com/tucnak/telebot"
    "github.com/mb-14/gomarkov"
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
)

const (
    tokensLengthLimit = 32
    order = 1
)

func main() {
    //Initialise
    //Get Telegram bot token from env
    bot_token := os.Getenv("TELEGRAM_BOT_TOKEN")

    //Initilise Telegram bot
    bot, err := telebot.NewBot(telebot.Settings{
        URL: "",
        Token: bot_token,
        Poller: &telebot.LongPoller{
            Timeout: 10 * time.Second,
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    //Initialise chain
    chain := gomarkov.NewChain(order)

    //Open source data dir
    //TODO Read from _all_ the files in the dir, not just data
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
    bot.Handle(telebot.OnText, func(m *telebot.Message) {
    //Process input

    //Identify whether to respond

    //Generate response
        sentence := generateSentence(chain, []string{"I", "think"})
        response := strings.Join(sentence, " ")

    //Respond with generated response
        bot.Send(m.Sender, response)
    })

    fmt.Println("Starting bot...")
    bot.Start()
    fmt.Println("Stopping bot...")
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
