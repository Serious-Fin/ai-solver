package problem

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"serious-fin/api/common"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetProblemsQueryThrowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewProblemHandler(db)

	mock.ExpectQuery("SELECT id, title FROM problems").WillReturnError(errors.New("error querying data"))

	if _, err = mockDb.GetProblems(); err == nil {
		t.Error("expected error when query fails")
	}

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

	var mockDb = NewProblemHandler(db)
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
		t.Errorf("unexpected error when returned rows are in a correct format: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

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

	var mockDb = NewProblemHandler(db)
	problemId := "1"
	want := &Problem{
		Id:          1,
		Title:       "foo",
		Description: "bar",
		TestCases: []common.TestCase{
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

	got, err := mockDb.GetProblemById(fmt.Sprint(want.Id))
	if err != nil {
		t.Errorf("unexpected error when returned rows are in a correct format: %v", err)
	}

	if res := reflect.DeepEqual(got, want); res == false {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMainFuncGo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var mockDb = NewProblemHandler(db)
	problemId := "1"
	want := "foo"

	values := [][]driver.Value{{want}}

	mock.ExpectQuery("SELECT mainFunction FROM goTemplates WHERE problemFk = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"mainFunction",
	}).AddRows(values...))

	got, err := mockDb.GetMainFuncGo(problemId)
	if err != nil {
		t.Errorf("unexpected error when returned rows are in a correct format: %v", err)
	}

	if want != got {
		t.Errorf("want: %v, got: %v", want, got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
