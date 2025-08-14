package user

import (
	"database/sql"
	"errors"
	"fmt"
	"serious-fin/api/common"
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	User User `json:"user"`
}

type SessionRequest struct {
	UserId string `json:"userId"`
}

type SessionResponse struct {
	SessionId string `json:"sessionId"`
}

type SessionInfoResponse struct {
	User User `json:"user"`
}

type User struct {
	Id         string `json:"id"`
	ProfilePic string `json:"profilePic"`
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
}

type Session struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
}

type UserDBHandler struct {
	DB common.DBInterface
}

func NewUserHandler(db common.DBInterface) *UserDBHandler {
	return &UserDBHandler{DB: db}
}

const sessionExpireDuration time.Duration = 7 * 24 * time.Hour

func (handler *UserDBHandler) GetUser(userId string) (*User, error) {
	row := handler.DB.QueryRow("SELECT id, email, name, profilePic FROM users WHERE id = ?", userId)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.ProfilePic)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not scan user (userId %s): %w", userId, err)
	}
	return &user, nil
}

func (handler *UserDBHandler) CreateUser(userInfo User) (*User, error) {
	row := handler.DB.QueryRow("INSERT INTO users (id, email, name, profilePic) VALUES (?, ?, ?, ?) RETURNING id, email, name, profilePic", userInfo.Id, userInfo.Email, userInfo.Name, userInfo.ProfilePic)
	var newUser User
	err := row.Scan(&newUser.Id, &newUser.Email, &newUser.Name, &newUser.ProfilePic)
	if err != nil {
		return nil, fmt.Errorf("could not insert new user: %w", err)
	}
	return &newUser, nil
}

func (handler *UserDBHandler) GetUserFromSession(sessionId string) (*User, error) {
	row := handler.DB.QueryRow("SELECT u.id, u.email, u.name, u.profilePic FROM sessions AS s LEFT JOIN users AS u ON s.userId = u.id WHERE s.id = ?", sessionId)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.ProfilePic)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not find user for session (sessionId %s): %w", sessionId, err)
	}
	return &user, nil
}

func (handler *UserDBHandler) GetSession(userId string) (*Session, error) {
	row := handler.DB.QueryRow("SELECT id, userId, expiresAt FROM sessions WHERE userId = ?", userId)
	var session Session
	err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not scan session (userId %s): %w", userId, err)
	}
	return &session, nil
}

func (handler *UserDBHandler) CreateSession(userId string) (*Session, error) {
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
		return nil, fmt.Errorf("could not insert session (userId %s): %w", userId, err)
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

func (handler *UserDBHandler) CleanupExpiredSessions(userId string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := handler.DB.Exec("DELETE FROM sessions WHERE userId = ? AND expiresAt < ?", userId, now)
	if err != nil {
		return fmt.Errorf("could not cleanup expired sessions: %w", err)
	}
	return nil
}
