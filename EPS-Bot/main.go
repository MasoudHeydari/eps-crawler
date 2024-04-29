package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	Start         = "/start"
	Location      = "📍 Location"
	EN            = "🇬🇧 En"
	GR            = "🇩🇪 Gr"
	FR            = "🇫🇷 Fr"
	ES            = "🇪🇸 Es" // Spanish
	PT            = "🇵🇹 Pt" // Portuguese
	RU            = "🇷🇺 Ru" // Portuguese
	Language      = "🌐 Language"
	GetLastResult = "📈 Get Last Results"
	APIToken      = "6704135678:AAH2SGETz7tSKY5NsQ-vv2-zO7tj-XZaKAk"
	TGBotPassword = "1234"
)

var mustProvidePassword = true

var (
	helpKeys = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GetLastResult)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Language),
			tgbotapi.NewKeyboardButton(Location),
		),
	)

	langKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(EN, "LangEN"),
			tgbotapi.NewInlineKeyboardButtonData(GR, "LangGR"),
			tgbotapi.NewInlineKeyboardButtonData(FR, "LangFR"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ES, "LangES"),
			tgbotapi.NewInlineKeyboardButtonData(PT, "LangPT"),
			tgbotapi.NewInlineKeyboardButtonData(RU, "LangRU"),
		),
	)

	locKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(EN, "LocEN"),
			tgbotapi.NewInlineKeyboardButtonData(GR, "LocGR"),
			tgbotapi.NewInlineKeyboardButtonData(FR, "LocFR"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ES, "LocES"),
			tgbotapi.NewInlineKeyboardButtonData(PT, "LocPT"),
			tgbotapi.NewInlineKeyboardButtonData(RU, "LocRU"),
		),
	)
)

func main() {
	bot, err := tgbotapi.NewBotAPI(APIToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if mustProvidePassword {
				if update.Message.Text == TGBotPassword {
					mustProvidePassword = false
					msg.ReplyMarkup = helpKeys
					msg.Text = "You are authenticated"
				} else {
					msg.Text = "Password is incorrect. Try again"
				}
			}

			// Construct a new message from the given chat ID and containing
			// the text that we received.
			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Text {
			case Start:
				if mustProvidePassword {
					msg.Text = "Welcome to EPS Crawler Bot.\nTo use this bot, please provide your password:"
				} else {
					msg.Text = "You are already authenticated.\nUse the provided keys to work with EPS Bot."
				}
			case Location:
				msg.Text = "📍 Select Your Location: "
				msg.ReplyMarkup = locKeyboard
			case Language:
				msg.Text = "🌐 Select Your Language: "
				msg.ReplyMarkup = langKeyboard
			case GetLastResult:
				msg.Text = "get last results"
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				// TODO: Don't panic!
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
