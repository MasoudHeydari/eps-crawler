package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"simple-telegram-bot/model"
	"strconv"
	"strings"
)

const (
	Start          = "/start"
	Location       = "📍 Location"
	EN             = "🇬🇧 En"
	GR             = "🇩🇪 Gr"
	FR             = "🇫🇷 Fr"
	ES             = "🇪🇸 Es" // Spanish
	PT             = "🇵🇹 Pt" // Portuguese
	RU             = "🇷🇺 Ru" // Portuguese
	Language       = "🌐 Language"
	GetLastResult  = "📈 Get Last Results"
	SetSearchQuery = "🔎 Set Your Search Query"
	APIToken       = "7053131148:AAG3NtM0ZJFxHEiGRQoUXeKlhDLqfVj5x78" //"6704135678:AAH2SGETz7tSKY5NsQ-vv2-zO7tj-XZaKAk"
	TGBotPassword  = "1234"
)

var (
	lastInput           string
	mustProvidePassword = true
	SQID                = 0
)

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
			tgbotapi.NewInlineKeyboardButtonData(EN, "🌐 EN Language Selected"),
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
			tgbotapi.NewInlineKeyboardButtonData(GR, "📍 GR Location Selected"),
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
	client := &http.Client{}
	sq := model.SearchQueryRequest{}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			messageText := update.Message.Text
			messages := make([]tgbotapi.MessageConfig, 0)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			messages = append(messages, msg)
			if mustProvidePassword {
				if update.Message.Text == TGBotPassword {
					mustProvidePassword = false
					msg.ReplyMarkup = helpKeys
					msg.Text = "You are authenticated"
				} else {
					msg.Text = "Password is incorrect. Try again"
				}
			}

			switch lastInput {
			case Start:
			case Location:
				sq.Location = messageText
			case Language:
				sq.Language = messageText
			case SetSearchQuery:
				sq.Query = messageText
			case GetLastResult:
				sqid, err := strconv.Atoi(messageText)
				if err != nil {
					msg.Text = "Search query ID must be a number"
					break
				}
				SQID = sqid
				if SQID == 0 {
					msg.Text = "Search query ID cannot be zero"
					break
				}
				msgs, err := GetLastResults(client, update.Message.Chat.ID, SQID)
				if err != nil {
					logrus.Info(err)
					break
				}
				messages = append(messages, msgs...)
			}

			// Construct a new message from the given chat ID and containing
			// the text that we received.
			// If the message was open, add a copy of our numeric keyboard.
			switch messageText {
			case Start:
				lastInput = Start
				if mustProvidePassword {
					msg.Text = "Welcome to EPS Crawler Bot.\nTo use this bot, please provide your password:"
				} else {
					msg.Text = "You are already authenticated.\nUse the provided keys to work with EPS Bot."
				}
			case Location:
				lastInput = Language
				msg.Text = "📍 Select Your Location: "
				msg.ReplyMarkup = locKeyboard
			case Language:
				lastInput = Language
				msg.Text = "🌐 Select Your Language: "
				msg.ReplyMarkup = langKeyboard
			case SetSearchQuery:
				lastInput = SetSearchQuery
				msg.Text = "Please provide your search query to let the crawler starts it's work:"
			case GetLastResult:
				lastInput = GetLastResult
				msg.Text = "Please provide your search query ID:"
			}

			// Send the message.
			for _, m := range messages {
				if _, err = bot.Send(m); err != nil {
					logrus.Info("failed to send message: %w", err)
				}
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

func GetLastResults(client *http.Client, msgChatID int64, sqID int) ([]tgbotapi.MessageConfig, error) {
	urlPath := fmt.Sprintf("http://localhost:9999/api/v1/search/%d", sqID)
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("GetLastResults.NewRequest: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		//msg.Text = "We have problem to connect to our servers, Please try again later"
		return nil, fmt.Errorf("GetLastResults.Get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetLastResults.StatusCode: unintented status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var serps []model.SERPResponse
	err = json.Unmarshal(bodyBytes, &serps)
	if err != nil {
		return nil, fmt.Errorf("GetLastResults.Unmarshal: %w", err)
	}
	tgMessages := make([]tgbotapi.MessageConfig, 0, len(serps))
	for _, serp := range serps {
		msg := tgbotapi.NewMessage(msgChatID, "")
		phones, emails, keyWords := "", "", ""
		if len(serp.ContactInfo.Phones) == 0 {
			phones = "Not Fount"
		} else {
			phones = strings.Join(serp.ContactInfo.Phones, ", ")
		}
		if len(serp.ContactInfo.Emails) == 0 {
			emails = "Not Fount"
		} else {
			emails = strings.Join(serp.ContactInfo.Emails, ", ")
		}
		if len(serp.Keywords) == 0 {
			keyWords = "Not Fount"
		} else {
			keyWords = strings.Join(serp.Keywords, ", ")
		}
		msg.Text = fmt.Sprintf(
			`%s

%s

%s

Key Words: %s

Emails: %s

Phones: %s
`,
			serp.URL,
			serp.Title,
			serp.Description,
			keyWords,
			emails,
			phones,
		)
		tgMessages = append(tgMessages, msg)
	}
	return tgMessages, nil
}
