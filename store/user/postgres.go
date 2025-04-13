package user

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

// NewPostgresStore creates an instance of UserStore over postgres connection.
func NewPostgresStore(pool *pgxpool.Pool) *postgresStore {
	return &postgresStore{
		pool: pool,
	}
}

const pgInsertOne = `
	insert into authgo.user(email, hash_password, username, first_name, last_name, middle_name)
	values ($1, $2, $3, $4, $5, $6)
	returning id;
`

func (s *postgresStore) InsertOne(ctx context.Context, user model.UserDao) (int64, error) {
	var userID int64

	err := s.pool.QueryRow(
		ctx,
		pgInsertOne,
		user.Email,
		user.HashPassword,
		user.Username,
		user.FirstName,
		user.LastName,
		user.MiddleName,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}

const pgFindOne = `
	select id, email, hash_password, username, first_name, last_name, middle_name, created_at, updated_at, is_deleted
	from authgo.user
	where id=$1;
`

func (s *postgresStore) FindOne(ctx context.Context, userID int64) (model.UserDao, error) {
	var user model.UserDao

	if err := s.pool.QueryRow(ctx, pgFindOne).Scan(&user); err != nil {
		return user, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

const pgUpdateOne = `
	update authgo.user
	set email=$2,
		hash_password=$3,
		username=$4,
		first_name=$5,
		last_name=$6,
		middle_name=$7
	where id=$1;
`

func (s *postgresStore) UpdateOne(ctx context.Context, user model.UserDao) error {
	res, err := s.pool.Exec(
		ctx,
		pgUpdateOne,
		user.ID,
		user.Email,
		user.HashPassword,
		user.Username,
		user.FirstName,
		user.LastName,
		user.MiddleName,
	)
	if err != nil {
		return fmt.Errorf("update user data: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("update user data: %w", pgx.ErrNoRows)
	}

	return nil
}

const pgDeleteOne = `
	delete from authgo.user
	where id=$1;
`

func (s *postgresStore) DeleteOne(ctx context.Context, userID int64) error {
	res, err := s.pool.Exec(
		ctx,
		pgDeleteOne,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("delete user: %w", pgx.ErrNoRows)
	}

	return nil
}

const pgListAll = `
	select id, email, hash_password, username, first_name, last_name, middle_name, created_at, updated_at, is_deleted
	from authgo.user;
`

func (s *postgresStore) ListAll(ctx context.Context) ([]model.UserDao, error) {
	rows, err := s.pool.Query(ctx, pgListAll)
	if err != nil {
		return nil, fmt.Errorf("list all users: %w", err)
	}

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.UserDao, error) {
		var user model.UserDao
		err := row.Scan(&user)
		return user, err
	})
	if err != nil {
		return nil, fmt.Errorf("list all users: %w", err)
	}

	return users, nil
}

const pgSetRole = `
	insert into authgo.user_role (user_id, role_id)
	values ($1, $2);
`

func (s *postgresStore) SetRole(ctx context.Context, userID, roleID int64) error {
	if _, err := s.pool.Exec(ctx, pgSetRole, userID, roleID); err != nil {
		return fmt.Errorf("insert user role: %w", err)
	}

	return nil
}

func (s *postgresStore) ApplyMigrations() error {
	if err := db.ApplyMigrations("postgres", db.PgMigrations, stdlib.OpenDBFromPool(s.pool)); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
