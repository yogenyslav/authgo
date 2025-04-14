package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/yogenyslav/authgo/db"
	"github.com/yogenyslav/authgo/model"
)

type roleStore struct {
	pg *postgresDB
}

// NewRoleStore creates an instance of RoleStore over postgres connection.
func NewRoleStore(pg *postgresDB) *roleStore {
	return &roleStore{
		pg: pg,
	}
}

func (s *roleStore) StartTx(ctx context.Context) (context.Context, error) {
	return s.pg.StartTx(ctx)
}

func (s *roleStore) CommitTx(ctx context.Context) error {
	return s.pg.CommitTx(ctx)
}

func (s *roleStore) RollbackTx(ctx context.Context) error {
	return s.pg.RollbackTx(ctx)
}

func (s *roleStore) ApplyMigrations() error {
	if err := db.ApplyMigrations("postgres", db.PgMigrations, stdlib.OpenDBFromPool(s.pg.GetPool())); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

const insertOneRole = `
	insert into authgo.role(name)
	values ($1)
	returning id;
`

func (s *roleStore) InsertOne(ctx context.Context, name string) (int64, error) {
	var roleID int64

	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return 0, fmt.Errorf("get conn: %w", err)
	}

	if err := conn.QueryRow(ctx, insertOneRole, name).Scan(&roleID); err != nil {
		return 0, fmt.Errorf("insert role: %w", err)
	}

	return roleID, nil
}

const findOneRoleByID = `
	select id, name, created_at
	from authgo.role
	where id=$1;
`

func (s *roleStore) FindOneByID(ctx context.Context, roleID int64) (model.RoleDao, error) {
	var role model.RoleDao

	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return role, fmt.Errorf("get conn: %w", err)
	}

	if err := conn.QueryRow(ctx, findOneRoleByID, roleID).Scan(&role); err != nil {
		return role, fmt.Errorf("find role: %w", err)
	}

	return role, nil
}

const findOneRoleByName = `
	select id, name, created_at
	from authgo.role
	where name=$1;
`

func (s *roleStore) FindOneByName(ctx context.Context, name string) (model.RoleDao, error) {
	var role model.RoleDao

	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return role, fmt.Errorf("get conn: %w", err)
	}

	if err := conn.QueryRow(ctx, findOneRoleByName, name).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
	); err != nil {
		return role, fmt.Errorf("find role: %w", err)
	}

	return role, nil
}

const updateOneRole = `
	update authgo.role
	set name=$2
	where id=$1;
`

func (s *roleStore) UpdateOne(ctx context.Context, role model.RoleDao) error {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn: %w", err)
	}

	res, err := conn.Exec(
		ctx,
		updateOneRole,
		role.ID,
		role.Name,
	)
	if err != nil {
		return fmt.Errorf("update role data: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("update role data: %w", pgx.ErrNoRows)
	}

	return nil
}

const deleteOneRole = `
	delete from authgo.role
	where id=$1;
`

func (s *roleStore) DeleteOne(ctx context.Context, roleID int64) error {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn: %w", err)
	}

	res, err := conn.Exec(
		ctx,
		deleteOneRole,
		roleID,
	)
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("delete role: %w", pgx.ErrNoRows)
	}

	return nil
}

const listAllRoles = `
	select id, name, created_at
	from authgo.role;
`

func (s *roleStore) ListAll(ctx context.Context) ([]model.RoleDao, error) {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn: %w", err)
	}

	rows, err := conn.Query(ctx, listAllRoles)
	if err != nil {
		return nil, fmt.Errorf("list all roles: %w", err)
	}

	roles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.RoleDao, error) {
		var role model.RoleDao
		err := row.Scan(&role.ID, &role.Name, &role.CreatedAt)
		return role, err
	})
	if err != nil {
		return nil, fmt.Errorf("list all roles: %w", err)
	}

	return roles, nil
}

const listUserRoles = `
	select r.id, r.name, r.created_at from authgo.role r
	join authgo.user_role ur
		on ur.role_id = r.id
	where ur.user_id = $1;
`

func (s *roleStore) ListUserRoles(ctx context.Context, userID int64) ([]model.RoleDao, error) {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn: %w", err)
	}

	rows, err := conn.Query(ctx, listUserRoles, userID)
	if err != nil {
		return nil, fmt.Errorf("list user roles: %w", err)
	}

	roles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.RoleDao, error) {
		var role model.RoleDao
		err := row.Scan(&role.ID, &role.Name, &role.CreatedAt)
		return role, err
	})
	if err != nil {
		return nil, fmt.Errorf("collect roles: %w", err)
	}

	return roles, nil
}
