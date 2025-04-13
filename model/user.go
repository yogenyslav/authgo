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

// UserRegister is a model of a Register request.
type UserRegister struct {
	Email      string
	Password   string
	Username   string
	FirstName  string
	LastName   string
	MiddleName string
}

// UserLogin is a model of a Login request.
type UserLogin struct {
	Email    string
	Password string
}

// AuthMeta is a model with data used to validate user's identity and permissions during requests.
type AuthMeta struct {
	UserID int64     `json:"sub"`
	Roles  []RoleDto `json:"roles"`
}

// AuthResp is a general response model for requests Login and Register.
type AuthResp struct {
	Token string
	Type  string
	Meta  AuthMeta
}
