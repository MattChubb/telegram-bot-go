package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/mb-14/gomarkov"
	"github.com/tucnak/telebot.v2"
    "github.com/TwinProduction/go-away"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	//Read params from command line
	chainFilePath := flag.String("chainfile", "", "Saved JSON chain file")
	chattiness := flag.Float64("chattiness", 0.1, "Chattiness (0-1, how often to respond unprompted)")
	debug := flag.Bool("debug", false, "Debug logging")
	order := flag.Int("order", 1, "Markov chain order. Use caution with values above 2")
	saveEvery := flag.Int("saveevery", 100, "Save every N messages")
	sourceDir := flag.String("sourcedir", "", "Source directory for training data")
	tokensLengthLimit := flag.Int("lengthlimit", 32, "Limit response length")
    trainOnly := flag.Bool("train-only", false, "Training only mode, do not attempt to connect to Telegram")
	flag.Parse()

	//Initialise
	rand.Seed(time.Now().Unix())
    if *debug {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.InfoLevel)
    }

	//Initialise chain
	log.Info("Initialising chain...")
	chain := gomarkov.NewChain(*order)
	if len(*chainFilePath) > 0 {
		log.Info("Loading chain from: ", *chainFilePath)
		chainFile, err := ioutil.ReadFile(*chainFilePath)
		if err != nil {
			log.Fatal(err)
		}
		chain.UnmarshalJSON(chainFile)
	}

	//Train
	//TODO Allow specifying a list of files instead of a directory
	if len(*sourceDir) > 0 {
		log.Info("Opening source data...")
		source_files, err := ioutil.ReadDir(*sourceDir)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Training chain on source data...")
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
    if *trainOnly {
        os.Exit(0)
    }

	//Initilise Telegram bot
	log.Info("Initialising bot...")
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

	//Connect Markov to Telegram
	//TODO Decouple the Markov implementation from the Telegram bot, allowing other techniques to be swapped in later
	log.Info("Adding chain to bot...")
	mNumber := 0
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		log.Debug("Received message: " + m.Text)
		parsedMessage := processString(m.Text)

		//Train on input (Ensures we always have a response for new words)
		chain.Add(parsedMessage)

		//Respond with generated response
		respond := decideWhetherToRespond(m, *chattiness, "@"+bot.Me.Username)

		if respond {
			log.Debug("Responding...")
			response := generateResponse(chain, parsedMessage, *tokensLengthLimit)
			log.Debug("Sending response: " + response)
			bot.Send(m.Chat, response)
		} else {
			log.Debug("Not responding")
		}

		mNumber++
		if mNumber > *saveEvery && len(*chainFilePath) > 0 {
			mNumber = 0
			saveChain(chain, *chainFilePath)
		}
	})

	log.Info("Starting bot...")
	bot.Start()
	log.Info("Bot stopped")
}

func saveChain(chain *gomarkov.Chain, file string) {
	log.Info("Saving chain...")
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
	return strings.Split(strings.ToLower(rawString), " ")
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
		subject = extractSubject(message)
	}
	//TODO Bi-directional generation using both a forwards and a backwards trained Markov chains
	//TODO Any other clever Markov hacks?
	sentence := generateSentence(chain, subject, lengthLimit)
    sentence[0] = strings.Title(sentence[0])
    return strings.Join(sentence, " ")
}

func generateSentence(chain *gomarkov.Chain, init []string, lengthLimit int) []string {
	// This function has been separated from response generation to allow bidirectional generation later

    //Train on the initial tokens to avoid unknown n-grams
    chain.Add(init)

    // The length of our initialisation chain needs to match the Markov order
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

        //TODO Implement a replacement wordfilter instead of just removing profanity
		if len(next) > 0 && ! goaway.IsProfane(next) {
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

func extractSubject(message []string) []string {
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
        if ! isStopWord(word) && word[0] != '@' {
            trimmedMessage = append(trimmedMessage, word)
        }
    }
    return trimmedMessage
}

func isStopWord(word string) bool {
    stopWords := []string{"the", "and", "to", "a", "i", "in", "be", "of", "that", "have", "it", }
    for _, stopWord := range stopWords {
        if word == stopWord {
            return true
        }
    }

    return false
}

func decideWhetherToRespond(m *telebot.Message, chattiness float64, name string) bool {
	respond := rand.Float64() < chattiness
	if !respond && m.Chat.Type == telebot.ChatPrivate {
		respond = true
		log.Debug("Respond: TRUE, private chat")
	} else if !respond {
		log.Debug("Respond: Not feeling chatty, checking for direct mention")
		for _, entity := range m.Entities {
			log.Debug("Respond: Found entity " + string(entity.Type))
			if entity.Type == telebot.EntityMention {
				mention := m.Text[entity.Offset : entity.Offset+entity.Length]
				log.Debug("Respond: Entity is " + mention)
				if mention == name {
					respond = true
					log.Debug("Respond: TRUE, @mentioned directly")
				}
			}
		}
		if !respond {
			log.Debug("Respond: FALSE, not feeling chatty")
		}
	} else {
		log.Debug("Respond: TRUE, feeling chatty")
	}

	return respond
}
