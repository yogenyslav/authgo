-- +goose Up
-- +goose StatementBegin
create schema authgo;
create table authgo.user (
	id bigserial primary key,
	email text unique not null,
	hash_password text not null,
	username text unique not null,
	first_name text not null default '',
	last_name text not null default '',
	middle_name text not null default '',
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	is_deleted bool not null default false
);
create index user_email on authgo.user using hash(email);
create index user_username on authgo.user using hash(username);

create or replace function authgo.update_user()
	returns trigger as
$BODY$
begin
	new.updated_at := current_timestamp;
return new;
end;
$BODY$
	language plpgsql;

create trigger trg_update_user
before update on authgo.user
for each row
execute function authgo.update_user();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema authgo;
drop table authgo.user;
drop function authgo.update_user;
drop trigger trg_update_user;
-- +goose StatementEnd
