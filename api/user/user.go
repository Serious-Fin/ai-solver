package user

import (
	"database/sql"
	"errors"
	"fmt"
	"serious-fin/api/common"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type Session struct {
	Id        string `json:"id"`
	UserId    int    `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
}

type UserDBHandler struct {
	DB common.DBInterface
}

func NewUserHandler(db common.DBInterface) *UserDBHandler {
	return &UserDBHandler{DB: db}
}

const sessionExpireDuration time.Duration = 7 * 24 * time.Hour

/*
	TODO: to test (IGNORE THIS TODO ITEM)

- If no rows error, -1 with NO error
- If other error, -1 with error
- If return id, that id gets returned with NO error
*/
func (handler *UserDBHandler) GetUser(email string) (*User, error) {
	row := handler.DB.QueryRow("SELECT id, email FROM users WHERE email = ?", email)
	var user User
	err := row.Scan(&user.Id, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not scan user (email %s): %w", email, err)
	}
	return &user, nil
}

/*
	TODO: to test (IGNORE THIS TODO ITEM)

- if inserting returns error, throw error
- if insert good -> no err
*/
func (handler *UserDBHandler) CreateUser(email string) (*User, error) {
	row := handler.DB.QueryRow("INSERT INTO users (email) VALUES (?) RETURNING id, email", email)
	var newUser User
	err := row.Scan(&newUser.Id, &newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("could not insert user (email %s): %w", email, err)
	}
	return &newUser, nil
}

/*
	TODO: to test (IGNORE THIS TODO ITEM)

- If no rows error, -1 with NO error
- If other error, -1 with error
- If return id, that id gets returned with NO error
*/
func (handler *UserDBHandler) GetSession(userId int) (*Session, error) {
	row := handler.DB.QueryRow("SELECT id, userId, expiresAt FROM sessions WHERE userId = ?", userId)
	var session Session
	err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not scan session (userId %d): %w", userId, err)
	}
	return &session, nil
}

/*
	TODO: to test (IGNORE THIS TODO ITEM)

- if inserting returns error, throw error
- if insert good -> no err
*/
func (handler *UserDBHandler) CreateSession(userId int) (*Session, error) {
	sessionId := uuid.New().String()
	expiresAt := time.Now().Add(sessionExpireDuration)
	expiresAtStr := expiresAt.Format(time.RFC3339)

	row := handler.DB.QueryRow(
		"INSERT INTO sessions (id, userId, expiresAt) VALUES (?, ?, ?) RETURNING id, userId, expiresAt",
		sessionId, userId, expiresAtStr,
	)

	var session Session
	err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("could not insert session (userId %d): %w", userId, err)
	}
	return &session, nil
}

/*
TODO: test
- if no session with userID found -> error
- if some error while updating -> error
- if all good return -> good
*/
func (handler *UserDBHandler) UpdateSession(userId int) (*Session, error) {
	expiresAt := time.Now().Add(sessionExpireDuration)
	expiresAtStr := expiresAt.Format(time.RFC3339)

	row := handler.DB.QueryRow(
		"UPDATE sessions SET expiresAt = ? WHERE userId = ? RETURNING id, userId, expiresAt",
		expiresAtStr, userId,
	)

	var session Session
	err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no session found for userId %d", userId)
		}
		return nil, fmt.Errorf("could not update session (userId %d): %w", userId, err)
	}
	return &session, nil
}

/*
TODO: test
- if bad format -> true
- if good format but more than limit -> expired
- if format bad and less than limit -> good
*/
func (handler *UserDBHandler) IsSessionExpired(session *Session) bool {
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(expiresAt)
}

/*
TODO: test:
- if deleting throws -> error
- if nothing -> good
*/
func (handler *UserDBHandler) CleanupExpiredSessions(userId int) error {
	now := time.Now().Format(time.RFC3339)
	_, err := handler.DB.Exec("DELETE FROM sessions WHERE userId = ? AND expiresAt < ?", userId, now)
	if err != nil {
		return fmt.Errorf("could not cleanup expired sessions: %w", err)
	}
	return nil
}

// TODO: change error wrapping to use %w (IGNORE THIS TODO ITEM)
// TODO: write tests for new endpoints (IGNORE THIS TODO ITEM)
