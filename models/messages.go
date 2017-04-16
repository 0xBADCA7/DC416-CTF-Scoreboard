package models

import (
	"database/sql"
	"time"
)

// Message contains messages left by admins for participants.
type Message struct {
	Id        int       `json:"-"`
	Content   string    `json:"message"`
	CreatedAt time.Time `json:"created"`
}

// NewMessage constructs a new message with an invalid id until it's saved.
func NewMessage(content string) Message {
	return Message{
		-1,
		content,
		time.Now(),
	}
}

// AllMessages obtains all of the messages left by admins, ordered by most recently posted.
func AllMessages(db *sql.DB) ([]Message, error) {
	messages := []Message{}
	rows, err := db.Query(QAllMessages)
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		m := Message{}
		err = rows.Scan(&m.Id, &m.Content, &m.CreatedAt)
		if err != nil {
			return []Message{}, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// DeleteAllMessages removes all messages left by admins from the database.
func DeleteAllMessages(db *sql.DB) error {
	_, err := db.Exec(QDeleteAllMessages)
	return err
}

// Save stores a new message to display to CTF participants.
func (m *Message) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveMessage, m.Content, m.CreatedAt)
	return err
}
