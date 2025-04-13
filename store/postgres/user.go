package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/yogenyslav/authgo/db"
	"github.com/yogenyslav/authgo/model"
)

type userStore struct {
	pg *postgresDB
}

// NewUserStore creates an instance of UserStore over postgres connection.
func NewUserStore(pg *postgresDB) *userStore {
	return &userStore{
		pg: pg,
	}
}

const insertOneUser = `
	insert into authgo.user(email, hash_password, username, first_name, last_name, middle_name)
	values ($1, $2, $3, $4, $5, $6)
	returning id;
`

func (s *userStore) InsertOne(ctx context.Context, user model.UserDao) (int64, error) {
	var userID int64

	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return 0, fmt.Errorf("get conn: %w", err)
	}

	err = conn.QueryRow(
		ctx,
		insertOneUser,
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

const findOneUserByEmail = `
	select id, email, hash_password, username, first_name, last_name, middle_name, created_at, updated_at, is_deleted
	from authgo.user
	where email=$1;
`

func (s *userStore) FindOneByEmail(ctx context.Context, email string) (model.UserDao, error) {
	var user model.UserDao

	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return user, fmt.Errorf("get conn: %w", err)
	}

	if err := conn.QueryRow(ctx, findOneUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.HashPassword,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsDeleted,
	); err != nil {
		return user, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

const updateOneUser = `
	update authgo.user
	set email=$2,
		hash_password=$3,
		username=$4,
		first_name=$5,
		last_name=$6,
		middle_name=$7
	where id=$1;
`

func (s *userStore) UpdateOne(ctx context.Context, user model.UserDao) error {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn: %w", err)
	}

	res, err := conn.Exec(
		ctx,
		updateOneUser,
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

const deleteOneUser = `
	delete from authgo.user
	where id=$1;
`

func (s *userStore) DeleteOne(ctx context.Context, userID int64) error {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn: %w", err)
	}

	res, err := conn.Exec(
		ctx,
		deleteOneUser,
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

const listAllUsers = `
	select id, email, hash_password, username, first_name, last_name, middle_name, created_at, updated_at, is_deleted
	from authgo.user;
`

func (s *userStore) ListAll(ctx context.Context) ([]model.UserDao, error) {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get conn: %w", err)
	}

	rows, err := conn.Query(ctx, listAllUsers)
	if err != nil {
		return nil, fmt.Errorf("list all users: %w", err)
	}

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.UserDao, error) {
		var user model.UserDao
		err := row.Scan(
			&user.ID,
			&user.Email,
			&user.HashPassword,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.MiddleName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsDeleted,
		)
		return user, err
	})
	if err != nil {
		return nil, fmt.Errorf("list all users: %w", err)
	}

	return users, nil
}

const setRole = `
	insert into authgo.user_role (user_id, role_id)
	values ($1, $2);
`

func (s *userStore) SetRole(ctx context.Context, userID, roleID int64) error {
	conn, err := s.pg.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("get conn: %w", err)
	}

	if _, err := conn.Exec(ctx, setRole, userID, roleID); err != nil {
		return fmt.Errorf("insert user role: %w", err)
	}

	return nil
}

func (s *userStore) StartTx(ctx context.Context) (context.Context, error) {
	return s.pg.StartTx(ctx)
}

func (s *userStore) CommitTx(ctx context.Context) error {
	return s.pg.CommitTx(ctx)
}

func (s *userStore) RollbackTx(ctx context.Context) error {
	return s.pg.RollbackTx(ctx)
}

func (s *userStore) ApplyMigrations() error {
	if err := db.ApplyMigrations("postgres", db.PgMigrations, stdlib.OpenDBFromPool(s.pg.GetPool())); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
