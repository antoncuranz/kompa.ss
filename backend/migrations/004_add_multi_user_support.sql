-- +goose Up
alter table "user"
add column jwt_sub uuid unique;

alter table trip
add column owner_id integer not null default 1 references "user" on delete cascade;
alter table trip
alter column owner_id drop default;

create table permissions
(
    user_id        integer not null references "user" on delete cascade,
    trip_id        integer not null references trip on delete cascade,
    read_sensitive boolean not null default false,
    write          boolean not null default false,
    primary key(user_id, trip_id)
);

-- +goose Down
DROP TABLE IF EXISTS permissions;
