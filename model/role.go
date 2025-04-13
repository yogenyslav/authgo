package model

import "time"

const (
	DefaultRole string = "default"
)

type RoleDao struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *RoleDao) ToDto() RoleDto {
	return RoleDto{
		ID:   r.ID,
		Name: r.Name,
	}
}

type RoleDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
