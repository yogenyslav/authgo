package role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/yogenyslav/authgo/db"
	"github.com/yogenyslav/authgo/model"
)

type postgresStore struct {
	pool *pgxpool.Pool
}

// NewPostgresStore creates an instance of RoleStore over postgres connection.
func NewPostgresStore(pool *pgxpool.Pool) *postgresStore {
	return &postgresStore{
		pool: pool,
	}
}

const pgInsertOne = `
	insert into authgo.role(name)
	values ($1)
	returning id;
`

func (s *postgresStore) InsertOne(ctx context.Context, name string) (int64, error) {
	var roleID int64

	err := s.pool.QueryRow(ctx, pgInsertOne, name).Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf("insert role: %w", err)
	}

	return roleID, nil
}

const pgFindOne = `
	select id, name, created_at
	from authgo.role
	where id=$1;
`

func (s *postgresStore) FindOne(ctx context.Context, roleID int64) (model.RoleDao, error) {
	var role model.RoleDao

	if err := s.pool.QueryRow(ctx, pgFindOne).Scan(&role); err != nil {
		return role, fmt.Errorf("find role: %w", err)
	}

	return role, nil
}

const pgUpdateOne = `
	update authgo.role
	set name=$2
	where id=$1;
`

func (s *postgresStore) UpdateOne(ctx context.Context, role model.RoleDao) error {
	res, err := s.pool.Exec(
		ctx,
		pgUpdateOne,
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

const pgDeleteOne = `
	delete from authgo.role
	where id=$1;
`

func (s *postgresStore) DeleteOne(ctx context.Context, roleID int64) error {
	res, err := s.pool.Exec(
		ctx,
		pgDeleteOne,
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

const pgListAll = `
	select id, name, created_at
	from authgo.role;
`

func (s *postgresStore) ListAll(ctx context.Context) ([]model.RoleDao, error) {
	rows, err := s.pool.Query(ctx, pgListAll)
	if err != nil {
		return nil, fmt.Errorf("list all roles: %w", err)
	}

	roles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.RoleDao, error) {
		var role model.RoleDao
		err := row.Scan(&role)
		return role, err
	})
	if err != nil {
		return nil, fmt.Errorf("list all roles: %w", err)
	}

	return roles, nil
}

func (s *postgresStore) ApplyMigrations() error {
	if err := db.ApplyMigrations("postgres", db.PgMigrations, stdlib.OpenDBFromPool(s.pool)); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
