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

// MessageModel is implemented by types that can store and modify Messages.
type MessageModel interface {
	Save(*Message) error
	All() ([]Message, error)
	Delete(*Message) error
}

// MessageModelDB implements MessageModel to handle Messages in a sqlite database.
type MessageModelDB struct {
	db *sql.DB
}

// NewMessage constructs a new message with an invalid id until it's saved.
func NewMessage(content string) Message {
	return Message{
		-1,
		content,
		time.Now(),
	}
}

// NewMessageModelDB constructs a new MessageModelDB with a database connection.
func NewMessageModelDB(db *sql.DB) MessageModelDB {
	return MessageModelDB{db}
}

// All obtains all of the messages left by admins, ordered by most recently posted.
func (self MessageModelDB) All() ([]Message, error) {
	messages := []Message{}
	rows, err := self.db.Query(QAllMessages)
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

// Delete removes a message left by admins from the database.
func (self MessageModelDB) Delete(message *Message) error {
	_, err := self.db.Exec(QDeleteMessage, message.Id)
	message.Id = -1
	return err
}

// Save stores a new message to display to CTF participants.
func (self MessageModelDB) Save(message *Message) error {
	_, err := self.db.Exec(QSaveMessage, message.Content, message.CreatedAt)
	if err != nil {
		return err
	}
	err = self.db.QueryRow(QLastInsertedId).Scan(&message.Id)
	return err
}
