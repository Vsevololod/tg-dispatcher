package domain

import "encoding/json"

// Update представляет корневой объект
type Update struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
	UUID     string
}

// Message содержит информацию о сообщении
type Message struct {
	MessageID       int64           `json:"message_id"`
	From            User            `json:"from"`
	Chat            Chat            `json:"chat"`
	Date            int64           `json:"date"`
	ReplyToMessage  *Message        `json:"reply_to_message,omitempty"`
	Text            string          `json:"text,omitempty"`
	Entities        []Entity        `json:"entities,omitempty"`
	LinkPreviewOpts *PreviewOptions `json:"link_preview_options,omitempty"`
}

// User содержит информацию о пользователе
type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat содержит информацию о чате
type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

// Entity представляет структуру для сущностей в сообщении
type Entity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

// PreviewOptions содержит настройки предпросмотра ссылки
type PreviewOptions struct {
	URL              string `json:"url"`
	PreferLargeMedia bool   `json:"prefer_large_media"`
}

// ParseUpdate парсит JSON в структуру Update
func ParseUpdate(jsonData []byte) (Update, error) {
	var update Update
	err := json.Unmarshal(jsonData, &update)
	if err != nil {
		return Update{}, err
	}
	return update, nil
}
