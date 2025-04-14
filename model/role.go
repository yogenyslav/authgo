package model

import "time"

const (
	// DefaultRole is a role that is assigned to every user on register.
	DefaultRole string = "default"
)

// RoleDat is a data model for role.
type RoleDao struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// ToDto converts a role data model into logical model for role.
func (r *RoleDao) ToDto() RoleDto {
	return RoleDto{
		ID:   r.ID,
		Name: r.Name,
	}
}

// RoleDto is logical model for role.
type RoleDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
