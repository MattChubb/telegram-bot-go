package main

import (
	"bufio"
	"encoding/json"
	"flag"
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

var debug bool

func main() {
	//Read params from command line
	chainFilePath := flag.String("chainfile", "", "Saved JSON chain file")
	chattiness := flag.Float64("chattiness", 0.1, "Chattiness (0-1, how often to respond unprompted)")
	flag.BoolVar(&debug, "debug", false, "Debug logging")
	order := flag.Int("order", 1, "Markov chain order. Use caution with values above 2")
	saveEvery := flag.Int("saveevery", 100, "Save every N messages")
	sourceDir := flag.String("sourcedir", "", "Source directory for training data")
	tokensLengthLimit := flag.Int("lengthlimit", 32, "Limit response length")
	flag.Parse()

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
	chain := gomarkov.NewChain(*order)
	if len(*chainFilePath) > 0 {
		fmt.Println("Loading chain from: ", *chainFilePath)
		chainFile, err := ioutil.ReadFile(*chainFilePath)
		if err != nil {
			log.Fatal(err)
		}
		chain.UnmarshalJSON(chainFile)
	}

	//Train
	//TODO Allow running in training-only mode for training models
	//TODO Allow specifying a list of files instead of a directory
	if len(*sourceDir) > 0 {
		fmt.Println("Opening source data...")
		source_files, err := ioutil.ReadDir(*sourceDir)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Training chain on source data...")
		for _, fileInfo := range source_files {
			if fileInfo.Name()[1] == '.' {
				continue
			}
			sourceFile, err := os.Open(*sourceDir + "/" + fileInfo.Name())
			if err != nil {
				log.Fatal(err)
			}
			trainFromFile(chain, sourceFile)
		}

		if len(*chainFilePath) > 0 {
			saveChain(chain, *chainFilePath)
		}
	}

	//Connect Markov to Telegram
	//TODO Decouple the Markov implementation from the Telegram bot, allowing other techniques to be swapped in later
	fmt.Println("Adding chain to bot...")
	mNumber := 0
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		logDebug("Received message: " + m.Text)
		parsedMessage := processString(m.Text)

		//Train on input (Ensures we always have a response for new words)
		chain.Add(parsedMessage)

		//Respond with generated response
        respond := decideWhetherToRespond(m, *chattiness, "@"+bot.Me.Username)

		if respond {
			logDebug("Responding...")
			response := generateResponse(chain, parsedMessage, *tokensLengthLimit)
			logDebug("Sending response: " + response)
			bot.Send(m.Chat, response)
		} else {
			logDebug("Not responding")
		}

		mNumber++
		if mNumber > *saveEvery && len(*chainFilePath) > 0 {
			mNumber = 0
			saveChain(chain, *chainFilePath)
		}
	})

	fmt.Println("Starting bot...")
	bot.Start()
	fmt.Println("Bot stopped")
}

func logDebug(e string) {
	if debug {
		fmt.Println(e)
	}
}

func saveChain(chain *gomarkov.Chain, file string) {
	fmt.Println("Saving chain...")
	chainJSON, err := json.Marshal(chain)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(file, chainJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func generateResponse(chain *gomarkov.Chain, message []string, lengthLimit int) string {
	subject := []string{}
	if len(message) > 0 {
		//TODO Do something cleverer with subject extraction
		subject = append([]string{}, message[rand.Intn(len(message))])
	}
	//TODO Bi-directional generation using both a forwards and a backwards trained Markov chains
	//TODO Any other clever Markov hacks?
	sentence := generateSentence(chain, subject, lengthLimit)
	return strings.Join(sentence, " ")
}

func generateSentence(chain *gomarkov.Chain, init []string, lengthLimit int) []string {
	// This function has been separated from response generation to allow bidirectional generation later
	tokens := []string{}
	if len(init) < chain.Order {
		for i := 0; i < chain.Order; i++ {
			tokens = append(tokens, gomarkov.StartToken)
		}
		tokens = append(tokens, init...)
	} else if len(init) > chain.Order {
		tokens = init[:chain.Order]
	} else {
		tokens = init
	}

	for tokens[len(tokens)-1] != gomarkov.EndToken &&
		len(tokens) < lengthLimit {
		next, err := chain.Generate(tokens[(len(tokens) - 1):])
		if err != nil {
			log.Fatal(err)
		}

		if len(next) > 0 {
			//TODO Add a wordfilter
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

func decideWhetherToRespond(m *telebot.Message, chattiness float64, name string) bool {
    respond := rand.Float64() < chattiness
    if !respond && m.Chat.Type == telebot.ChatPrivate {
        respond = true
        logDebug("Respond: TRUE, private chat")
    } else if !respond {
        logDebug("Respond: Not feeling chatty, checking for direct mention")
        for _, entity := range m.Entities {
            logDebug("Respond: Found entity " + string(entity.Type))
            if entity.Type == telebot.EntityMention {
                mention := m.Text[entity.Offset : entity.Offset+entity.Length]
                logDebug("Respond: Entity is " + mention)
                if mention == name {
                    respond = true
                    logDebug("Respond: TRUE, @mentioned directly")
                }
            }
        }
        if !respond {
            logDebug("Respond: FALSE, not feeling chatty")
        }
    } else {
        logDebug("Respond: TRUE, feeling chatty")
    }

    return respond
}
