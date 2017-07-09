package models

import (
	"database/sql"
	"fmt"
	"time"
)

// SessionCookieName is the name to give the cookie stored in the user's browser
// containing their session token.
const SessionCookieName = "session"

// Admin session tokens will be allowed to live for 12 hours.
const tokenLifetime = 12 * time.Hour

// Session contains information about a session token created when someone logs
// into the admin dashboard.
type Session struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

// SessionModel is implemented by types that can create and invalidate sessions.
type SessionModel interface {
	Find(string) (Session, error)
	Save(*Session) error
	Delete(*Session) error
}

// SessionModelDB implements SessionModel to handle Sessions in a sqlite database.
type SessionModelDB struct {
	db *sql.DB
}

// NewSession creates a new Session instance with a created time and expiration time, but
// without an actual token. Token generation occurs when the token is saved, since the
// operation may result in an error.
func NewSession() Session {
	now := time.Now()
	return Session{
		Token:   "",
		Created: now,
		Expires: now.Add(tokenLifetime),
	}
}

// NewSessionModelDB constructs a new instance of SessionModelDB with a database connection.
func NewSessionModelDB(db *sql.DB) SessionModelDB {
	return SessionModelDB{db}
}

// Find attempts to lookup a session given its unique token.
func (self SessionModelDB) Find(token string) (Session, error) {
	session := Session{}
	fmt.Println("Looking for session token", token)
	err := self.db.QueryRow(QFindSessionToken, token).Scan(&session.Created, &session.Expires)
	if err != nil {
		return Session{}, err
	}
	session.Token = token
	return session, err
}

// Save attempts to save a session token, first generating the actual token value itself,
// guaranteeing it is unique.
func (self SessionModelDB) Save(session *Session) error {
	uniqueToken := generateUniqueToken(func(token string) bool {
		_, err := self.Find(token)
		return err != nil
	})
	session.Token = uniqueToken
	_, err := self.db.Exec(QCreateSession, session.Token, session.Created, session.Expires)
	return err
}

// Delete removes a session token from the database.
func (self SessionModelDB) Delete(session *Session) error {
	_, err := self.db.Exec(QDeleteSession, session.Token)
	session.Token = ""
	return err
}

// IsExpired determines if a session token is expired and should, consequently, be
// rejected for administrative applicaions.
func (s Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}
