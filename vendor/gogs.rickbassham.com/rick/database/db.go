package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	// ErrUnregisteredStatement is returned when there is an attempt to execute a statement that
	// was not previously registered.
	ErrUnregisteredStatement = errors.New("statement not registered")
)

// DB represents the functions needed to access the db. This is satisfied by the sqlx.DB struct.
type DB interface {
	execer
	selecter
	getter
	Ping() error
	Preparex(query string) (*sqlx.Stmt, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// Tx represents the functions needed to operate in a transaction. This is satisfied by the sqlx.Tx struct.
type Tx interface {
	execer
	selecter
	getter
	Commit() error
	Rollback() error
}

// Middleware represents something you want to run on every query. This could be to start a
// newrelic segment, log the query, or anything else you want.
type Middleware interface {
	Before(ctx context.Context, name, statement string, args ...interface{}) (context.Context, error)
	After(ctx context.Context, err error, name, statement string, args ...interface{}) error
}

// Database is our wrapper around sqlx which enforces correctness of queries.
type Database struct {
	statements map[string]string
	middleware []Middleware
	db         DB
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type selecter interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type getter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// New creates a new wrapper around a sqlx.DB.
func New(db DB) (*Database, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{
		statements: map[string]string{},
		db:         db,
	}, nil
}

// With will add the middleware to the Database wrapper.
func (db *Database) With(mw ...Middleware) *Database {
	for _, m := range mw {
		db.middleware = append(db.middleware, m)
	}

	return db
}

// RegisterStatement registers the given statement with the given name. It validates the statement
// by using the Preparex func.
func (db *Database) RegisterStatement(name, statement string) error {
	err := db.ValidateStatement(statement)
	if err != nil {
		return err
	}

	db.statements[name] = statement

	return nil
}

// ValidateStatement is used to verify the syntax of the given statement.
func (db *Database) ValidateStatement(statement string) error {
	prepared, err := db.db.Preparex(statement)
	if err != nil {
		return err
	}

	// mocks will return a nil prepared statement, but all real implementations will return an error if prepared is nil.
	if prepared != nil {
		err = prepared.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) exec(ctx context.Context, e execer, statementName string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.getStatement(statementName)
	if err != nil {
		return nil, err
	}

	for _, mw := range db.middleware {
		ctx, err = mw.Before(ctx, statementName, stmt, args...)
		if err != nil {
			return nil, err
		}
	}

	if len(args) > 0 {
		stmt, args, err = sqlx.In(stmt, args...)
		if err != nil {
			return nil, err
		}
	}

	res, err := e.ExecContext(ctx, stmt, args...)

	for _, mw := range db.middleware {
		err = mw.After(ctx, err, statementName, stmt, args...)
		if err != nil {
			return nil, err
		}
	}

	return res, err
}

// BeginTx will start a db transaction with the given options.
func (db *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return db.db.BeginTxx(ctx, opts)
}

// Insert is used to insert data into the database. It retuns the id of the inserted record.
func (db *Database) Insert(ctx context.Context, statementName string, args ...interface{}) (lastInsertID int64, err error) {
	result, err := db.exec(ctx, db.db, statementName, args...)
	if err != nil {
		return 0, err
	}

	lastInsertID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// InsertTx is used to insert data into the database inside a transaction. It returns the id of the inserted record.
func (db *Database) InsertTx(ctx context.Context, tx Tx, statementName string, args ...interface{}) (lastInsertID int64, err error) {
	result, err := db.exec(ctx, tx, statementName, args...)
	if err != nil {
		return 0, err
	}

	lastInsertID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// Update updates records in the database. It returns the number of rows affected by the update.
func (db *Database) Update(ctx context.Context, statementName string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := db.exec(ctx, db.db, statementName, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// UpdateTx updates records in the database inside a transaction. It returns the number of rows affected by the update.
func (db *Database) UpdateTx(ctx context.Context, tx Tx, statementName string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := db.exec(ctx, tx, statementName, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Delete deletes records from the database. It returns the number of rows affected by the delete.
func (db *Database) Delete(ctx context.Context, statementName string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := db.exec(ctx, db.db, statementName, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// DeleteTx deletes records from the database inside a transaction. It returns the number of rows affected by the delete.
func (db *Database) DeleteTx(ctx context.Context, tx Tx, statementName string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := db.exec(ctx, tx, statementName, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Exec executes an arbitrary statement on the database. It returns a sql.Result.
func (db *Database) Exec(ctx context.Context, statementName string, args ...interface{}) (sql.Result, error) {
	return db.exec(ctx, db.db, statementName, args...)
}

// ExecTx executes an arbitrary statement on the database in a transaction. It returns a sql.Result.
func (db *Database) ExecTx(ctx context.Context, tx Tx, statementName string, args ...interface{}) (sql.Result, error) {
	return db.exec(ctx, tx, statementName, args...)
}

func (db *Database) execSelect(ctx context.Context, s selecter, dest interface{}, statementName string, args ...interface{}) error {
	stmt, err := db.getStatement(statementName)
	if err != nil {
		return err
	}

	for _, mw := range db.middleware {
		ctx, err = mw.Before(ctx, statementName, stmt, args...)
		if err != nil {
			return err
		}
	}

	if len(args) > 0 {
		stmt, args, err = sqlx.In(stmt, args...)
		if err != nil {
			return err
		}
	}

	err = s.SelectContext(ctx, dest, stmt, args...)

	for _, mw := range db.middleware {
		err = mw.After(ctx, err, statementName, stmt, args...)
		if err != nil {
			return err
		}
	}

	return err
}

// Select will run a select query on the database and bind the results to dest.
func (db *Database) Select(ctx context.Context, dest interface{}, statementName string, args ...interface{}) error {
	return db.execSelect(ctx, db.db, dest, statementName, args...)
}

// SelectTx will run a select query on the database inside a transaction and bind the results to dest.
func (db *Database) SelectTx(ctx context.Context, tx Tx, dest interface{}, statementName string, args ...interface{}) error {
	return db.execSelect(ctx, tx, dest, statementName, args...)
}

func (db *Database) execGet(ctx context.Context, g getter, dest interface{}, statementName string, args ...interface{}) error {
	stmt, err := db.getStatement(statementName)
	if err != nil {
		return err
	}

	for _, mw := range db.middleware {
		ctx, err = mw.Before(ctx, statementName, stmt, args...)
		if err != nil {
			return err
		}
	}

	if len(args) > 0 {
		stmt, args, err = sqlx.In(stmt, args...)
		if err != nil {
			return err
		}
	}

	err = g.GetContext(ctx, dest, stmt, args...)

	for _, mw := range db.middleware {
		err = mw.After(ctx, err, statementName, stmt, args...)
		if err != nil {
			return err
		}
	}

	return err
}

// Get will run a select query to get 1 record from the database and bind the result to dest.
func (db *Database) Get(ctx context.Context, dest interface{}, statementName string, args ...interface{}) error {
	return db.execGet(ctx, db.db, dest, statementName, args...)
}

// GetTx will run a select query to get 1 record from the database in a transaction and bind the result to dest.
func (db *Database) GetTx(ctx context.Context, tx Tx, dest interface{}, statementName string, args ...interface{}) error {
	return db.execGet(ctx, tx, dest, statementName, args...)
}

func (db *Database) getStatement(statementName string) (string, error) {
	stmt, ok := db.statements[statementName]
	if !ok {
		return "", ErrUnregisteredStatement
	}

	_, file, line, _ := runtime.Caller(3)

	idx := strings.LastIndexByte(file, '/')
	if idx != -1 {
		idx = strings.LastIndexByte(file[:idx], '/')
	}

	if idx != -1 {
		file = file[idx+1:]
	}

	stmt = fmt.Sprintf("%s /* %s:%d */", stmt, file, line)

	return stmt, nil
}
