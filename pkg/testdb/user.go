package testdb

import (
	"context"

	"github.com/rickbassham/example-go/pkg/identity"
)

type User struct {
	Entity

	Username string `db:"username"`
}

func init() {
	statements["user_insert"] = "INSERT INTO users (created_by, updated_by, username) VALUES (?, ?, ?)"
	statements["user_select_active"] = "SELECT id, created_at, updated_at, created_by, updated_by, username FROM users WHERE deleted_at IS NULL"
}

func (db *DB) InsertUser(ctx context.Context, username string) (int, error) {
	execUser := identity.FromContext(ctx)

	id, err := db.db.Insert(ctx, "user_insert", execUser, execUser, username)

	return int(id), err
}

func (db *DB) GetActiveUsers(ctx context.Context) ([]User, error) {
	var users []User

	err := db.db.Select(ctx, &users, "user_select_active")
	if err != nil {
		return nil, err
	}

	return users, nil
}
