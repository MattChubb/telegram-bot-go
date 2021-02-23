# telegram-bot-go
![Test](https://github.com/MattChubb/telegram-bot-go/workflows/Go/badge.svg?branch=master&event=push)
![CodeQL](https://github.com/MattChubb/telegram-bot-go/workflows/CodeQL/badge.svg?branch=master&event=push)

Markov chain based Telegram bot

# Usage
Telegram credentials are passed through via the TELEGRAM_BOT_TOKEN environment variable. All other options are passed in via the command line and are set up with "sensible defaults" that will enable you to get up and running quickly and adjust from there. To see these command line options in full, use:
    telegrambot -help

Any source data should be in text files, separated by newlines. An empty "source_data" directory is provided in this repo for convenience. The source data (if specified) will be used to train the bot on startup.

The bot will also learn from every message it receives (even if it doesn't reply to it). Over time, this means it will begin to imitate the conversation styles of any chats it participates in. While it will never explicitly store messages it recieves verbatim, it is not recommended to place it anywhere it might overhear sensitive details.

For help setting up a Telegram bot, check their guide here: https://core.telegram.org/bots

# Features
## Brain types
## Chattiness
## Profanity filter
