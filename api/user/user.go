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

func (handler *UserDBHandler) CreateUser(email string) (*User, error) {
	row := handler.DB.QueryRow("INSERT INTO users (email) VALUES (?) RETURNING id, email", email)
	var newUser User
	err := row.Scan(&newUser.Id, &newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("could not insert user (email %s): %w", email, err)
	}
	return &newUser, nil
}

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

func (handler *UserDBHandler) UpdateSession(sessionId string) (*Session, error) {
	expiresAt := time.Now().Add(sessionExpireDuration)
	expiresAtStr := expiresAt.Format(time.RFC3339)

	row := handler.DB.QueryRow(
		"UPDATE sessions SET expiresAt = ? WHERE id = ? RETURNING id, userId, expiresAt",
		expiresAtStr, sessionId,
	)

	var session Session
	err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no session found (session id %s)", sessionId)
		}
		return nil, fmt.Errorf("could not update session (session id %s): %w", sessionId, err)
	}
	return &session, nil
}

func IsSessionExpired(session *Session) bool {
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(expiresAt)
}

func (handler *UserDBHandler) CleanupExpiredSessions(userId int) error {
	now := time.Now().Format(time.RFC3339)
	_, err := handler.DB.Exec("DELETE FROM sessions WHERE userId = ? AND expiresAt < ?", userId, now)
	if err != nil {
		return fmt.Errorf("could not cleanup expired sessions: %w", err)
	}
	return nil
}
