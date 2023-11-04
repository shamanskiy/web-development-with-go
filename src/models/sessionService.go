package models

import (
	"database/sql"
	"fmt"
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB           *sql.DB
	TokenManager TokenManager
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, err := ss.TokenManager.New()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.TokenManager.Hash(token),
	}

	row := ss.DB.QueryRow(`
	  INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE SET token_hash = $2 
		RETURNING id;`,
		session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.TokenManager.Hash(token)
	var user User
	row := ss.DB.QueryRow(`
	  SELECT u.id, u.email, u.password_hash
		FROM users u JOIN sessions s ON u.id = s.user_id
		WHERE s.token_hash = $1;`,
		tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.TokenManager.Hash(token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
