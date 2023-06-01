package events

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"telegram-bot/client"
	"telegram-bot/storage"
)

const (
	startCommand  = "/start"
	helpCommand   = "/help"
	addCommand    = "/add"
	pickCommand   = "pick_random"
	religionCat   = "Religion"
	philosophyCat = "Philosophy"
	workCat       = "Work"
	familyCat     = "Family"
)

var buttonsMain = [][]client.KeyboardButton{
	{client.KeyboardButton{Text: addCommand}, client.KeyboardButton{Text: "pick_random"}}}

var buttonsCat = [][]client.KeyboardButton{
	{client.KeyboardButton{Text: religionCat}, client.KeyboardButton{Text: philosophyCat}},
	{client.KeyboardButton{Text: workCat}, client.KeyboardButton{Text: familyCat}}}

func (h *Handler) processCommand(text, username string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("recieved new command: %s from %s", text, username)

	if category, ok := h.SelectedCategory[chatID]; ok {
		if h.CommandType == addCommand {
			delete(h.SelectedCategory, chatID)
			return h.addQuote(chatID, category, text, username)
		} else {
			delete(h.SelectedCategory, chatID)
			return h.pickQuote(chatID, category, username)
		}
	}
	switch text {
	case startCommand:
		return h.TgClient.SendMessage(chatID, ScriptStart, client.NewReplyKeyboardMarkup(buttonsMain), client.NewReplyKeyboardRemove(false))
	case helpCommand:
		return h.TgClient.SendMessage(chatID, ScriptHelp, client.NewReplyKeyboardMarkup(buttonsMain), client.NewReplyKeyboardRemove(false))
	case addCommand:
		return h.handleCommand(chatID, addCommand)
	case pickCommand:
		return h.handleCommand(chatID, pickCommand)
	case philosophyCat:
		return h.handleCategory(chatID, philosophyCat)
	case religionCat:
		return h.handleCategory(chatID, religionCat)
	case workCat:
		return h.handleCategory(chatID, workCat)
	case familyCat:
		return h.handleCategory(chatID, familyCat)
	default:
		return h.TgClient.SendMessage(chatID, ScriptUnknownCommand, client.NewReplyKeyboardMarkup(nil), client.NewReplyKeyboardRemove(true))
	}
}

func (h *Handler) handleCommand(chatID int, commandType string) error {
	if commandType == addCommand {
		h.CommandType = addCommand
		return h.TgClient.SendMessage(chatID, ScriptAdd, client.NewReplyKeyboardMarkup(buttonsCat), client.NewReplyKeyboardRemove(false))
	} else {
		h.CommandType = pickCommand
		return h.TgClient.SendMessage(chatID, ScriptPick, client.NewReplyKeyboardMarkup(buttonsCat), client.NewReplyKeyboardRemove(false))
	}
}

func (h *Handler) handleCategory(chatID int, category string) error {
	h.SelectedCategory[chatID] = category
	return h.TgClient.SendMessage(chatID, fmt.Sprintf(ScriptHandleCat, category), client.NewReplyKeyboardMarkup(nil), client.NewReplyKeyboardRemove(true))
}

func (h *Handler) addQuote(chatID int, category, text, username string) error {

	quote := &storage.Quote{
		Text:       text,
		CategoryId: indexCategories(category),
		Username:   username,
	}
	isExist, err := h.Storage.IfExists(context.Background(), quote)
	if err != nil {
		return err
	}
	if isExist {
		return h.TgClient.SendMessage(chatID, ScriptAlreadyAdded, client.NewReplyKeyboardMarkup(buttonsMain), client.NewReplyKeyboardRemove(false))
	}
	if err := h.Storage.Add(context.Background(), quote); err != nil {
		return err
	}
	return h.TgClient.SendMessage(chatID, ScriptAddingFinished, client.NewReplyKeyboardMarkup(buttonsMain), client.NewReplyKeyboardRemove(false))
}

func (h *Handler) pickQuote(chatID int, category, username string) error {
	quote, err := h.Storage.PickRandom(context.Background(), indexCategories(category), username)
	if err != nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return h.TgClient.SendMessage(chatID, ScriptQuoteNotFount, client.NewReplyKeyboardMarkup(buttonsMain), client.NewReplyKeyboardRemove(false))
	}
	return h.TgClient.SendMessage(chatID, quote.Text, client.NewReplyKeyboardMarkup(buttonsMain), client.NewReplyKeyboardRemove(false))
}

func indexCategories(category string) int {
	if category == religionCat {
		return 1
	} else if category == philosophyCat {
		return 2
	} else if category == workCat {
		return 3
	} else {
		return 4
	}
}
