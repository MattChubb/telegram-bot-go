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
	//TODO Turn these into command line params
    chattiness        = 0.5
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
	//TODO Allow loading a saved chain instead of training a new one

	//Train
	//TODO Allow running in training-only mode for training models
	//TODO Allow skipping of training step (helpful if we plan to load a saved chain)
	fmt.Println("Opening source data...")
	//TODO Allow specifying a list of files via command line
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
	//TODO Save chain to json file

	//Connect Markov to Telegram
	//TODO Decouple the Markov implementation from the Telegram bot, allowing other techniques to be swapped in later
	fmt.Println("Adding chain to bot...")
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		parsedMessage := processString(m.Text)

		//Train on input (Ensures we always have a response for new words)
		chain.Add(parsedMessage)

		//Respond with generated response
        respond := rand.Float32() < chattiness
        if !respond && m.Chat.Type == telebot.ChatPrivate {
            respond = true
        } else {
            for _, entity := range m.Entities {
                if entity.Type == telebot.EntityTMention && entity.User.ID == bot.Me.ID {
                    respond = true
                }
            }
        }

        if respond {
		    response := generateResponse(chain, parsedMessage)
		    bot.Send(m.Chat, response)
        }
	})

	fmt.Println("Starting bot...")
	bot.Start()
	fmt.Println("Bot stopped")
}

func processString(rawString string) []string {
	//TODO Handle punctuation other than spaces
	//TODO Lowercase everything to reduce the token space
	return strings.Split(rawString, " ")
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

func generateResponse(chain *gomarkov.Chain, message []string) string {
	subject := []string{}
	if len(message) > 0 {
		//TODO Do something cleverer with subject extraction
		subject = append([]string{}, message[rand.Intn(len(message))])
	}
	//TODO Bi-directional generation using both a forwards and a backwards trained Markov chains
	//TODO Any other clever Markov hacks?
	sentence := generateSentence(chain, subject)
	return strings.Join(sentence, " ")
}

func generateSentence(chain *gomarkov.Chain, init []string) []string {
	// This function has been separated from response generation to allow bidirectional generation later
	tokens := []string{}
	if len(init) < chain.Order {
        for i:=0; i < chain.Order; i++ {
            tokens = append(tokens, gomarkov.StartToken)
        }
        tokens = append(tokens, init...)
	} else if len(init) > chain.Order {
        tokens = init[:chain.Order]
    } else {
        tokens = init
    }

	for tokens[len(tokens)-1] != gomarkov.EndToken &&
		len(tokens) < tokensLengthLimit {
		next, err := chain.Generate(tokens[(len(tokens) - 1):])
	if err != nil {
	    log.Fatal(err)
	}

		if len(next) > 0 {
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
