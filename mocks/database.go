package mocks

import (
	"context"
	"fmt"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/** Mock Implementation for Database Client to call Raw SQL Queries */

/**
 * This is a compile-time interface compliance check in Go.
 * Ensures MockDatabasePool implements repository.DBTX interface
 */
var _ repository.DBTX = (*MockDatabasePool)(nil)

type MockDatabasePool struct {
	ExecFunc     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func (m *MockDatabasePool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, args...)
	}

	return pgconn.CommandTag{}, nil
}

func (m *MockDatabasePool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}

	return nil
}

func (m *MockDatabasePool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}

	return nil, nil
}

/* MockRow implements pgx.Row for testing */
type MockRow struct {
	ScanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil
}

/* MockRows implements pgx.Rows for testing */
type MockRows struct {
	rows     [][]interface{}
	current  int
	closeErr error
}

func NewMockRows(rows [][]interface{}) *MockRows {
	return &MockRows{
		rows:    rows,
		current: -1,
	}
}

func (m *MockRows) Next() bool {
	m.current++
	return m.current < len(m.rows)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if m.current >= len(m.rows) {
		return fmt.Errorf("no more rows")
	}

	row := m.rows[m.current]
	for i, v := range row {
		if i >= len(dest) {
			break
		}

		/* Type assertion and assignment */
		switch d := dest[i].(type) {
		case *string:
			*d = v.(string)
		case *int64:
			*d = v.(int64)
		}
	}
	return nil
}

func (m *MockRows) Close() {}

func (m *MockRows) Err() error {
	return m.closeErr
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (m *MockRows) Values() ([]any, error) {
	return nil, nil
}

func (m *MockRows) RawValues() [][]byte {
	return nil
}

func (m *MockRows) Conn() *pgx.Conn {
	return nil
}
