package problem

import (
	"errors"
	"testing"
)

type RowMock struct {
	ScanFunc  func(dest ...any) error
	NextFunc  func() bool
	ErrFunc   func() error
	CloseFunc func() error
}

func (m RowMock) Scan(dest ...any) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil
}

type DbMock struct {
	QueryFunc    func(query string, args ...any) (*RowMock, error)
	QueryRowFunc func(query string, args ...any) *RowMock
}

func (m DbMock) Query(query string, args ...any) (*RowMock, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return &RowMock{}, nil
}

func (m DbMock) QueryRow(query string, args ...any) *RowMock {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return &RowMock{}
}

func TestQueryThrowsError(t *testing.T) {
	handler := NewProblemDBHandler(DbMock{
		QueryFunc: func(query string, args ...any) (*RowMock, error) {
			return nil, errors.New("error querying data")
		},
	})

	_, err := handler.GetProblems()
	if err == nil {
		t.Error("Expected to get error but received no error")
	}
}
