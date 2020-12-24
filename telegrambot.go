package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
    brain "github.com/MattChubb/telegram-bot-go/brain"
    doublemarkov "github.com/MattChubb/telegram-bot-go/brain/doublemarkov"
    //markov "github.com/MattChubb/telegram-bot-go/brain/markov"
)

func main() {
	//Read params from command line
	brainFilePath := flag.String("brainfile", "", "Saved JSON brain file")
	chattiness := flag.Float64("chattiness", 0.1, "Chattiness (0-1, how often to respond unprompted)")
	debug := flag.Bool("debug", false, "Debug logging")
	order := flag.Int("order", 2, "Markov brain order. Use caution with values above 2")
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

	//Initialise brain
	log.Info("Initialising brain...")
    //TODO Switch brain types via config
    brain := new(doublemarkov.Brain)
    brain.Init(*order, *tokensLengthLimit)
	if len(*brainFilePath) > 0 {
		log.Info("Loading brain from: ", *brainFilePath)
		brainFile, err := ioutil.ReadFile(*brainFilePath)
		if err != nil {
			log.Fatal(err)
		}
		brain.UnmarshalJSON(brainFile)
	}

	//Train
	//TODO Allow specifying a list of files instead of a directory
	if len(*sourceDir) > 0 {
        //TODO Split training from source files into its own method
		//TODO Add debug logging
		log.Info("Opening source data...")
		//TODO Is there any advantage to using ioutil over os.Readdir?
		source_files, err := ioutil.ReadDir(*sourceDir)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Training on source data...")
		for _, fileInfo := range source_files {
			if fileInfo.Name()[1] == '.' {
				continue
			}
			sourceFile, err := os.Open(*sourceDir + "/" + fileInfo.Name())
			if err != nil {
				log.Fatal(err)
			}
			trainFromFile(brain, sourceFile)
			//TODO Close file after reading!
		}

		if len(*brainFilePath) > 0 {
			saveBrain(brain, *brainFilePath)
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
	log.Info("Adding brain to bot...")
	mNumber := 0
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
        //TODO Split handler into its own method
		log.Debug("Received message: " + m.Text)

		//Train on input (Ensures we always have a response for new words)
		brain.Train(m.Text)

		//Respond with generated response
		respond := decideWhetherToRespond(m, *chattiness, "@"+bot.Me.Username)

		if respond {
			log.Debug("Responding...")
			response, _ := brain.Generate(m.Text)
			log.Debug("Sending response: " + response)
			bot.Send(m.Chat, response)
		} else {
			log.Debug("Not responding")
		}

		mNumber++
		if mNumber > *saveEvery && len(*brainFilePath) > 0 {
			mNumber = 0
			saveBrain(brain, *brainFilePath)
		}
	})

	log.Info("Starting bot...")
	bot.Start()
	log.Info("Bot stopped")
}

func saveBrain(brain brain.Brain, file string) {
	log.Info("Saving brain...")
	brainJSON, err := json.Marshal(brain)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(file, brainJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func trainFromFile(brain brain.Brain, file *os.File) {
	scanner := bufio.NewScanner(file)
    buffer := make([]byte, 0, 64*1024)
    scanner.Buffer(buffer, 1024*1024)
	for scanner.Scan() {
		brain.Train(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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
