package main

import (
    "github.com/tucnak/telebot"
    "github.com/mb-14/gomarkov"
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "os"
    "strings"
    "time"
)

const (
    tokensLengthLimit = 32
    order = 1
    sourceDir = "./source_data"
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

    //Misc init
    rand.Seed(time.Now().Unix())

    //Train
    source_files, err := ioutil.ReadDir(sourceDir)
    if err != nil {
        log.Fatal(err)
    }

    for _, fileInfo := range source_files {
        if fileInfo.Name()[1] == '.' {
            continue
        }
        sourceFile, err := os.Open(sourceDir + "/" + fileInfo.Name())
        if err != nil {
            log.Fatal(err)
        }
        trainFromFile(chain, sourceFile)
    }

    //Connect Markov to Telegram
    bot.Handle(telebot.OnText, func(m *telebot.Message) {
        //Process input
        rawMessage := m.Text
        parsedMessage := processString(rawMessage)

        //Train on input (Ensures we always have a response for new words)
        chain.Add(parsedMessage)

        //Identify whether to respond
        respondable := true
        if !respondable {
            return
        }

        //Identify subject
        subject := []string{}
        subject = append(subject, parsedMessage[rand.Intn(len(parsedMessage))])

        //Generate response
        sentence := generateSentence(chain, subject)
        response := strings.Join(sentence, " ")

        //Respond with generated response
        bot.Send(m.Sender, response)
    })

    fmt.Println("Starting bot")
    bot.Start()
    fmt.Println("Bot stopped")
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

func trainFromFile(chain *gomarkov.Chain, file *os.File) {
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        chain.Add(processString(scanner.Text()))
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
