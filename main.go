package main

import (
	"bufio"
	"fmt"
	"github.com/mb-14/gomarkov"
	"github.com/tucnak/telebot"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	tokensLengthLimit = 32
	order             = 1
	sourceDir         = "./source_data"
)

func main() {
	//Initialise
	rand.Seed(time.Now().Unix())

	//Initilise Telegram bot
	fmt.Println("Initialising bot...")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := telebot.NewBot(telebot.Settings{
		URL:   "",
		Token: botToken,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	//Initialise chain
	fmt.Println("Initialising chain...")
	chain := gomarkov.NewChain(order)

	//Train
	fmt.Println("Opening source data...")
	source_files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Training chain on source data...")
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
	fmt.Println("Adding chain to bot...")
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		parsedMessage := processString(m.Text)

		//Train on input (Ensures we always have a response for new words)
		chain.Add(parsedMessage)

		//TODO Only speak when spoken to
		respondable := true
		if !respondable {
			return
		}

		//Generate response
		subject := append([]string{}, parsedMessage[rand.Intn(len(parsedMessage))])
		sentence := generateSentence(chain, subject)
		response := strings.Join(sentence, " ")

		//Respond with generated response
		bot.Send(m.Sender, response)
	})

	fmt.Println("Starting bot...")
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

	for tokens[len(tokens)-1] != gomarkov.EndToken &&
		len(tokens) < tokensLengthLimit {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
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
