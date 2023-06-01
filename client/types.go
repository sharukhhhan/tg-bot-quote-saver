package client

type Updates struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []Updates `json:"result"`
}

type Message struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

func NewReplyKeyboardMarkup(keyboard [][]KeyboardButton) *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{Keyboard: keyboard}
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
}

func NewReplyKeyboardRemove(removeKeyboard bool) *ReplyKeyboardRemove {
	return &ReplyKeyboardRemove{RemoveKeyboard: removeKeyboard}
}
