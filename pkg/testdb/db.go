package testdb

import (
	"time"

	"gogs.rickbassham.com/rick/database"
)

var (
	statements map[string]string
)

type DB struct {
	db *database.Database
}

func New(db *database.Database) (*DB, error) {
	err := initializeStatements(db)
	if err != nil {
		return nil, err
	}

	db.With(Logger{})
	db.With(NewRelic{})

	return &DB{
		db: db,
	}, nil
}

func initializeStatements(db *database.Database) error {
	for k, v := range statements {
		err := db.RegisterStatement(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

type Entity struct {
	ID        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	CreatedBy string     `db:"created_by"`
	UpdatedBy string     `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
}
