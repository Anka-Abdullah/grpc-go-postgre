package testdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"sync/atomic"
)

var counter int32

type Stub struct {
	QueryFunc func(query string, args []driver.NamedValue) ([][]driver.Value, error)
	ExecFunc  func(query string, args []driver.NamedValue) (lastInsertID int64, rowsAffected int64, err error)
}

// New creates a new sql.DB backed by a stub driver.
func New() (*sql.DB, *Stub, error) {
	drv := &Stub{}
	name := fmt.Sprintf("stub-%d", atomic.AddInt32(&counter, 1))
	sql.Register(name, drv)
	db, err := sql.Open(name, "")
	if err != nil {
		return nil, nil, err
	}
	return db, drv, nil
}

func (s *Stub) Open(name string) (driver.Conn, error) { return &conn{s}, nil }

type conn struct{ stub *Stub }

func (c *conn) Prepare(query string) (driver.Stmt, error) { return nil, fmt.Errorf("not implemented") }
func (c *conn) Close() error                              { return nil }
func (c *conn) Begin() (driver.Tx, error)                 { return nil, fmt.Errorf("not implemented") }

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.stub.QueryFunc == nil {
		return &rows{}, nil
	}
	vals, err := c.stub.QueryFunc(query, args)
	return &rows{values: vals}, err
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.stub.ExecFunc == nil {
		return stubResult{}, nil
	}
	lid, ra, err := c.stub.ExecFunc(query, args)
	return stubResult{lid: lid, ra: ra}, err
}

// Implement required interfaces
var _ driver.Driver = (*Stub)(nil)
var _ driver.Conn = (*conn)(nil)
var _ driver.QueryerContext = (*conn)(nil)
var _ driver.ExecerContext = (*conn)(nil)

type rows struct {
	values [][]driver.Value
	idx    int
}

func (r *rows) Columns() []string { return []string{} }
func (r *rows) Close() error      { return nil }
func (r *rows) Next(dest []driver.Value) error {
	if r.idx >= len(r.values) {
		return io.EOF
	}
	row := r.values[r.idx]
	for i := range row {
		dest[i] = row[i]
	}
	return nil
}

type stubResult struct{ lid, ra int64 }

func (r stubResult) LastInsertId() (int64, error) { return r.lid, nil }
func (r stubResult) RowsAffected() (int64, error) { return r.ra, nil }
