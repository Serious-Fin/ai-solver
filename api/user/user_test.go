package user

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

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
