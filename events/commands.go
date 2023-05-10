package events

import (
	"log"
	"strings"
	"telegram-bot/client"
)

const (
	startCommand  = "/start"
	helpCommand   = "/help"
	addCommand    = "/add"
	religionCat   = "Religion"
	philosophyCat = "Philosophy"
	workCat       = "Work"
	familyCat     = "Family"
)

func (h *Handler) processCommand(text, username string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("recieved new command: %s from %s", text, username)

	switch text {
	case startCommand:
		buttons := [][]client.KeyboardButton{
			{client.KeyboardButton{Text: "Button 1"}, client.KeyboardButton{Text: "Button 2"}},
			{client.KeyboardButton{Text: "Button 3"}, client.KeyboardButton{Text: "Button 4"}},
		}
		markup := client.ReplyKeyboardMarkup{
			Keyboard: buttons,
		}
		return h.tgClient.SendMessage(chatID, ScriptStart, markup)
	case helpCommand:
		buttons := [][]client.KeyboardButton{
			{client.KeyboardButton{Text: "Button 1"}, client.KeyboardButton{Text: "Button 2"}},
			{client.KeyboardButton{Text: "Button 3"}, client.KeyboardButton{Text: "Button 4"}},
		}
		markup := client.ReplyKeyboardMarkup{
			Keyboard: buttons,
		}
		return h.tgClient.SendMessage(chatID, ScriptHelp, markup)
	case addCommand:
		return h.addCommand(chatID, username)
	default:
		return h.tgClient.SendMessage(chatID, ScriptUnknownCommand, client.ReplyKeyboardMarkup{})
	}
}

func (h *Handler) addCommand(chatID int, username string) error {

	buttons := [][]client.KeyboardButton{
		{client.KeyboardButton{Text: "Religion"}, client.KeyboardButton{Text: "Philosophy"}},
		{client.KeyboardButton{Text: "Work"}, client.KeyboardButton{Text: "Family"}},
	}
	markup := client.ReplyKeyboardMarkup{
		Keyboard: buttons,
	}
	return h.tgClient.SendMessage(chatID, ScriptAdd, markup)
}

func (h *Handler) chooseCategory(chatID int, username string, categoryId int) error {
	buttons := [][]client.KeyboardButton{
		{client.KeyboardButton{Text: "Male"}, client.KeyboardButton{Text: "Female"}},
	}
	markup := client.ReplyKeyboardMarkup{
		Keyboard: buttons,
	}
	return h.tgClient.SendMessage(chatID, ScriptAdd2, markup)
}
