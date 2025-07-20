package problem

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetProblemsQueryThrowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewProblemDBHandler(db)

	mock.ExpectQuery("SELECT id, title FROM problems").WillReturnError(errors.New("error querying data"))

	if _, err = mockDb.GetProblems(); err == nil {
		t.Error("expected error when query fails")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProblems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewProblemDBHandler(db)
	want := []Problem{
		{
			Id:    1,
			Title: "one",
		},
		{
			Id:    2,
			Title: "two",
		},
	}

	values := [][]driver.Value{
		{
			want[0].Id, want[0].Title,
		},
		{
			want[1].Id, want[1].Title,
		},
	}

	mock.ExpectQuery("SELECT id, title FROM problems").WillReturnRows(sqlmock.NewRows([]string{
		"id", "title",
	}).AddRows(values...))

	got, err := mockDb.GetProblems()
	if err != nil {
		t.Error("unexpected error when returned rows are in a correct format")
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProblemById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewProblemDBHandler(db)
	problemId := "1"
	want := Problem{
		Id:          1,
		Title:       "foo",
		Description: "bar",
		TestCases: []TestCase{
			{
				Id: 0,
				Inputs: []string{
					"[]int{2, 7, 11, 15}", "9",
				},
				ExpectedOutput: "[]int{0, 1}",
			},
		},
	}

	values := [][]driver.Value{
		{
			want.Id, want.Title, want.Description, `[{"id": 0,"inputs":  ["[]int{2, 7, 11, 15}","9"],"output": "[]int{0, 1}"}]`,
		},
	}

	mock.ExpectQuery("SELECT id, title, description, testCases FROM problems WHERE id = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"id", "title", "description", "testCases",
	}).AddRows(values...))

	got, err := mockDb.GetProblemById(string(rune(want.Id)))
	if err != nil {
		t.Error("unexpected error when returned rows are in a correct format")
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
