package models

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMessageModelDB(test *testing.T) {
	db, err := sql.Open("sqlite3", "testmessagemodel.db")
	if err != nil {
		test.Error(err)
	}
	defer os.Remove("testmessagemodel.db")
	if err = InitTables(db); err != nil {
		test.Error(err)
	}
	messageModel := NewMessageModelDB(db)

	testCases := []struct {
		Message string
	}{
		{"Hello, world!"},
		{"Testing the MessageModelDB type"},
		{""},
	}

	// Expect All() to return zero messages before any are inserted.
	messages, err := messageModel.All()
	if err != nil {
		test.Error(err)
	}
	if len(messages) != 0 {
		test.Errorf("Expected 0 messages before inserting any. Got %d\n", len(messages))
	}
	for _, testCase := range testCases {
		msg := NewMessage(testCase.Message)
		err := messageModel.Save(msg)
		if err != nil {
			test.Error(err)
		}
		if msg.Id < 0 {
			test.Errorf("Expected ID of newly-created message to be updated after save. It was not")
		}
		messages = append(messages, msg)
	}

	// Expect Delete() to remove items, and for All() to return less messages each time one is removed.
	for i := 0; i < len(messages); i++ {
		found, err := messageModel.All()
		if err != nil {
			test.Error(err)
		}
		if len(found) != len(messages) {
			test.Errorf("Expected %d messages. Got %d\n", len(messages), len(found))
		}
		err = messageModel.Delete(messages[i])
		if err != nil {
			test.Error(err)
		}
		if messages[i].Id > 0 {
			test.Errorf("Expected ID of deleted message to be invalidated")
		}
	}

	// Expect database operations to fail after closing the database connection.
	db.Close()
	msg := NewMessage("should fail")
	err = messageModel.Save(msg)
	if err == nil {
		test.Errorf("Expected to get an error when saving after closing the database connection, but we did not")
	}
	_, err = messageModel.All()
	if err == nil {
		test.Errorf("Expected to get an error retrieving messages after closing the database connection, but we did not")
	}
	err = messageModel.Delete(msg)
	if err == nil {
		test.Errorf("Expected to get an error deleting after closing the database connection, but we did not")
	}
}
