package models

import (
	"database/sql"
	"fmt"
	"time"
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

const (
	// DefaultResetDuration is the default time that a PasswordReset is
	// valid for.
	DefaultResetDuration = 1 * time.Hour
)

type PasswordResetService struct {
	DB           *sql.DB
	TokenManager TokenManager
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	token, err := prs.TokenManager.New()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	passwordReset := PasswordReset{
		Token:     token,
		TokenHash: prs.TokenManager.Hash(token),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	row := prs.DB.QueryRow(`
	  SELECT id 
	  FROM users WHERE email=$1`, email)
	err = row.Scan(&passwordReset.UserID)
	if err != nil {
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	row = prs.DB.QueryRow(`
	  INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE SET token_hash = $2, expires_at = $3 
		RETURNING id;`,
		passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)
	err = row.Scan(&passwordReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	return &passwordReset, nil
}

// We are going to consume a token and return the user associated with it, or return an error if the token wasn't valid for any reason.
func (prs *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := prs.TokenManager.Hash(token)
	var user User
	row := prs.DB.QueryRow(`
	  SELECT u.id, u.email, u.password_hash
		FROM users u JOIN password_resets pr ON u.id = pr.user_id
		WHERE pr.token_hash = $1 AND $2 < expires_at;`,
		tokenHash, time.Now())
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("comsume password reset: %w", err)
	}

	_, err = prs.DB.Exec(`
		DELETE FROM password_resets
		WHERE user_id = $1;`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("delete password reset: %w", err)
	}

	return &user, nil
}
