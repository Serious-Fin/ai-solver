package user

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserNoRowsReturnsNil(t *testing.T) {
	userEmail := "example@abc.com"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("SELECT id, email FROM users WHERE email = ?").WithArgs(userEmail).WillReturnRows(sqlmock.NewRows([]string{
		"id", "email",
	}))

	user, err := mockDb.GetUser(userEmail)
	if err != nil {
		t.Errorf("unexpected error when no user found: %v", err)
	}

	if user != nil {
		t.Errorf("user should be nil, but was %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserErrorThrowsError(t *testing.T) {
	userEmail := "example@abc.com"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("SELECT id, email FROM users WHERE email = ?").WithArgs(userEmail).WillReturnError(errors.New("something happened"))

	if _, err := mockDb.GetUser(userEmail); err == nil {
		t.Error("expected error when reading from db throws error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserReturnsUser(t *testing.T) {
	want := &User{
		Id:    1,
		Email: "example.abc.com",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("SELECT id, email FROM users WHERE email = ?").WithArgs(want.Email).WillReturnRows(sqlmock.NewRows([]string{
		"id", "email",
	}).AddRow([]driver.Value{want.Id, want.Email}...))

	got, err := mockDb.GetUser(want.Email)
	if err != nil {
		t.Errorf("unexpected error when reading from db does not throw: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUserErrorThrowsError(t *testing.T) {
	userEmail := "example@abc.com"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("INSERT INTO users .* VALUES .* RETURNING id, email").WithArgs(userEmail).WillReturnError(errors.New("something happened"))

	if _, err := mockDb.CreateUser(userEmail); err == nil {
		t.Error("expected error when executing from db throws error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUserReturnsNewUser(t *testing.T) {
	want := &User{
		Id:    1,
		Email: "example.abc.com",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("INSERT INTO users .* VALUES .* RETURNING id, email").WithArgs(want.Email).WillReturnRows(sqlmock.NewRows([]string{
		"id", "email",
	}).AddRow([]driver.Value{want.Id, want.Email}...))

	got, err := mockDb.CreateUser(want.Email)
	if err != nil {
		t.Errorf("unexpected error when reading from db does not throw: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSessionNoRowsReturnsNil(t *testing.T) {
	userId := 1
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("SELECT id, userId, expiresAt FROM sessions WHERE userId = ?").WithArgs(userId).WillReturnRows(sqlmock.NewRows([]string{
		"id", "email",
	}))

	user, err := mockDb.GetSession(userId)
	if err != nil {
		t.Errorf("unexpected error when no session found: %v", err)
	}

	if user != nil {
		t.Errorf("session should be nil, but was %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSessionErrorThrowsError(t *testing.T) {
	userId := 1
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("SELECT id, userId, expiresAt FROM sessions WHERE userId = ?").WithArgs(userId).WillReturnError(errors.New("something happened"))

	if _, err := mockDb.GetSession(userId); err == nil {
		t.Error("expected error when reading from db throws error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSessionReturnsSession(t *testing.T) {
	want := &Session{
		Id:        "1",
		UserId:    1,
		ExpiresAt: "2006-01-02T15:04:05Z07:00",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("SELECT id, userId, expiresAt FROM sessions WHERE userId = ?").WithArgs(want.UserId).WillReturnRows(sqlmock.NewRows([]string{
		"id", "userId", "expiresAt",
	}).AddRow([]driver.Value{want.Id, want.UserId, want.ExpiresAt}...))

	got, err := mockDb.GetSession(want.UserId)
	if err != nil {
		t.Errorf("unexpected error when reading from db does not throw: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateSessionErrorThrowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("INSERT INTO sessions .* VALUES .* RETURNING id, userId, expiresAt").WillReturnError(errors.New("something happened"))

	if _, err := mockDb.CreateSession(1); err == nil {
		t.Error("expected error when executing from db throws error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateSessionReturnsNewSession(t *testing.T) {
	want := &Session{
		Id:        "1",
		UserId:    1,
		ExpiresAt: "2006-01-02T15:04:05Z07:00",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("INSERT INTO sessions .* VALUES .* RETURNING id, userId, expiresAt").WillReturnRows(sqlmock.NewRows([]string{
		"id", "userId", "expiresAt",
	}).AddRow([]driver.Value{want.Id, want.UserId, want.ExpiresAt}...))

	got, err := mockDb.CreateSession(want.UserId)
	if err != nil {
		t.Errorf("unexpected error when reading from db does not throw: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIsSessionExpiredBadFormat(t *testing.T) {
	if expired := IsSessionExpired(&Session{
		ExpiresAt: "bad format",
	}); expired == false {
		t.Errorf("expected session with bad time format to be considered expired, but it was not")
	}
}

func TestIsSessionExpiredOldSession(t *testing.T) {
	if expired := IsSessionExpired(&Session{
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
	}); expired == true {
		t.Errorf("expected session to be considered expired, but it was not")
	}
}

func TestIsSessionExpiredActiveSession(t *testing.T) {
	if expired := IsSessionExpired(&Session{
		ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
	}); expired == false {
		t.Errorf("expected session to be considered still active, but it was not")
	}
}

func TestUpdateSessionNoSessionByIdError(t *testing.T) {
	sessionId := "sessionId"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("UPDATE sessions SET .* WHERE .* RETURNING id, userId, expiresAt").WillReturnRows(sqlmock.NewRows([]string{
		"id", "userId", "expiresAt",
	}))

	if _, err := mockDb.UpdateSession(sessionId); err == nil {
		t.Error("expected error when no session found, but got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateSessionErrorThrowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("UPDATE sessions SET .* WHERE .* RETURNING id, userId, expiresAt").WillReturnError(errors.New("something happened"))

	if _, err := mockDb.UpdateSession("sessionId"); err == nil {
		t.Error("expected error when executing from db throws error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateSessionReturnsSession(t *testing.T) {
	want := &Session{
		Id:        "1",
		UserId:    1,
		ExpiresAt: "2006-01-02T15:04:05Z07:00",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectQuery("UPDATE sessions SET .* WHERE .* RETURNING id, userId, expiresAt").WillReturnRows(sqlmock.NewRows([]string{
		"id", "userId", "expiresAt",
	}).AddRow([]driver.Value{want.Id, want.UserId, want.ExpiresAt}...))

	got, err := mockDb.UpdateSession(want.Id)
	if err != nil {
		t.Errorf("unexpected error when reading from db does not throw: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCleanupExpiredSessionsErrorThrowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectExec("DELETE FROM sessions WHERE .* AND .*").WillReturnError(errors.New("something happened"))

	if err := mockDb.CleanupExpiredSessions(1); err == nil {
		t.Error("expected error when executing from db throws error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCleanupExpiredSessionsSuccessReturnsNil(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewUserHandler(db)

	mock.ExpectExec("DELETE FROM sessions WHERE .* AND .*").WillReturnResult(sqlmock.NewResult(0, 1))

	if err := mockDb.CleanupExpiredSessions(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
