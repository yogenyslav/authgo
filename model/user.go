package model

import "time"

// UserDao is a user model in data store.
type UserDao struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	HashPassword string    `db:"hash_password"`
	Username     string    `db:"username"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	MiddleName   string    `db:"middle_name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	IsDeleted    bool      `db:"is_deleted"`
}

type UserRegister struct {
	Email      string
	Password   string
	Username   string
	FirstName  string
	LastName   string
	MiddleName string
}
