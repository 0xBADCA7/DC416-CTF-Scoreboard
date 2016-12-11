package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Admin session tokens will be allowed to live for 12 hours.
const tokenLifetime = 12 * time.Hour

// Session contains information about a session token created when someone logs
// into the admin dashboard.
type Session struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

// FindSession attempts to lookup a session given its unique token.
func FindSession(db *sql.DB, token string) (Session, error) {
	session := Session{}
	fmt.Println("Looking for session token", token)
	err := db.QueryRow(QFindSessionToken, token).Scan(&session.Created, &session.Expires)
	if err != nil {
		return Session{}, err
	}
	session.Token = token
	return session, err
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

// Save attempts to save a session token, first generating the actual token value itself,
// guaranteeing it is unique.
func (s *Session) Save(db *sql.DB) error {
	uniqueToken := generateUniqueToken(func(token string) bool {
		_, err := FindSession(db, token)
		return err != nil
	})
	s.Token = uniqueToken
	_, err := db.Exec(QCreateSession, s.Token, s.Created, s.Expires)
	return err
}

// Delete removes a session token from the database.
func (s *Session) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteSession, s.Token)
	return err
}

// IsExpired determines if a session token is expired and should, consequently, be
// rejected for administrative applicaions.
func (s *Session) IsExpired() bool {
	return s.Expires.After(time.Now())
}
