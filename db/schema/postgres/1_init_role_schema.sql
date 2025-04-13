-- +goose Up
-- +goose StatementBegin
create schema authgo;
create table authgo.role (
	id bigserial primary key,
	name text not null,
	created_at timestamp not null default current_timestamp,
	is_deleted bool not null default false
);
create index role_name on authgo.role using hash(name);

insert into authgo.role (name)
values ('default');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table authgo.role;
drop schema authgo;
-- +goose StatementEnd
