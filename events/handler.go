package events

import (
	"errors"
	"fmt"
	"telegram-bot/client"
)

type Handler struct {
	tgClient *client.Client
	offset   int
}

func NewHandler(tgClient *client.Client) *Handler {
	return &Handler{tgClient: tgClient}
}

func (h *Handler) Fetch(limit int) ([]Event, error) {
	updates, err := h.tgClient.Updates(h.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]Event, 0, len(updates))

	for _, upd := range updates {
		result = append(result, convertToEvent(upd))
	}

	h.offset = updates[len(updates)-1].ID + 1
	return result, nil
}

func (h *Handler) Process(event Event) error {
	switch event.Type {
	case Message:
		return h.processMessage(event)
	default:
		return fmt.Errorf("can't process the message :%w", errors.New("unknown event typ "))
	}
}

func (h *Handler) processMessage(event Event) error {
	err := h.processCommand(event.Text, event.Username, event.ChatId)
	if err != nil {
		return err
	}
	return nil
}

func convertToEvent(upd client.Updates) Event {
	var fetchedType Type
	var fetchedText string

	if upd.Message == nil {
		fetchedType = Unknown
		fetchedText = ""
	} else {
		fetchedType = Message
		fetchedText = upd.Message.Text
	}

	result := Event{
		Type: fetchedType,
		Text: fetchedText,
	}

	if fetchedType == Message {
		result.ChatId = upd.Message.Chat.ID
		result.Username = upd.Message.From.Username
	}

	return result

}
