package user

import (
	"database/sql"
	"errors"
	"fmt"
	"serious-fin/api/common"
	"time"
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

/*
NEEDED OPERATIONS:
(BOTH EXPOSED AS A SINGLE ENDPOINT TO GET A USER, WE CREATE OT INTERNALLY)
GET /user/get?email=EMAIL [NEED TESTING]
POST /user/create body: {email: EMAIL} [NEED TESTING]

GET /session?userId=USERID [MISSING]
UPDATE /session body: {expireAt: NOW + 1 WEEK} [MISSING]
DELETE /session?userId=USERID [MISSING]
POST /session {userId: USERID} [MISSING]
*/

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

// -----------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------

// CreateSession creates a new session for a user
// POST /session {userId: USERID}
func (handler *UserDBHandler) CreateSession(userId int) (*Session, error) {
	sessionId := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 week from now
	expiresAtStr := expiresAt.Format(time.RFC3339)

	row := handler.DB.QueryRow(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?) RETURNING id, user_id, expires_at",
		sessionId, userId, expiresAtStr,
	)

	var session Session
	err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("could not insert session (userId %d): %w", userId, err)
	}
	return &session, nil
}

// UpdateSession updates the expiration time of a session
// UPDATE /session body: {expireAt: NOW + 1 WEEK}
func (handler *UserDBHandler) UpdateSession(userId int) (*Session, error) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 week from now
	expiresAtStr := expiresAt.Format(time.RFC3339)

	row := handler.DB.QueryRow(
		"UPDATE sessions SET expires_at = ? WHERE user_id = ? RETURNING id, user_id, expires_at",
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

// DeleteSession removes a session by userId
// DELETE /session?userId=USERID
func (handler *UserDBHandler) DeleteSession(userId int) error {
	result, err := handler.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userId)
	if err != nil {
		return fmt.Errorf("could not delete session (userId %d): %w", userId, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected for delete session (userId %d): %w", userId, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no session found to delete for userId %d", userId)
	}

	return nil
}

// Helper method to check if a session is expired
func (handler *UserDBHandler) IsSessionExpired(session *Session) bool {
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return true // If we can't parse the time, consider it expired
	}
	return time.Now().After(expiresAt)
}

// CleanupExpiredSessions removes all expired sessions from the database
func (handler *UserDBHandler) CleanupExpiredSessions() (int64, error) {
	now := time.Now().Format(time.RFC3339)
	result, err := handler.DB.Exec("DELETE FROM sessions WHERE expires_at < ?", now)
	if err != nil {
		return 0, fmt.Errorf("could not cleanup expired sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get rows affected for cleanup: %w", err)
	}

	return rowsAffected, nil
}

// -----------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------------------------------------

// TODO: write endpoints for users and sessions (IGNORE THIS TODO ITEM)
// TODO: change error wrapping to use %w (IGNORE THIS TODO ITEM)
// TODO: write tests for new endpoints (IGNORE THIS TODO ITEM)
