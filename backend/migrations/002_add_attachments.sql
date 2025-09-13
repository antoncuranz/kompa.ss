-- +goose Up
create table attachment
(
    id                  serial primary key,
    trip_id             integer not null references trip on delete cascade,
    name                varchar(255) not null,
    blob                bytea not null
);

create table attachment_activity
(
    attachment_id       integer not null references attachment on delete cascade,
    activity_id         integer not null references activity on delete cascade
);

create table attachment_accommodation
(
    attachment_id       integer not null references attachment on delete cascade,
    accommodation_id    integer not null references accommodation on delete cascade
);

create table attachment_transportation
(
    attachment_id       integer not null references attachment on delete cascade,
    transportation_id   integer not null references transportation on delete cascade
);

-- +goose Down
DROP TABLE IF EXISTS attachment;
DROP TABLE IF EXISTS attachment_activity;
DROP TABLE IF EXISTS attachment_accommodation;
DROP TABLE IF EXISTS attachment_transportation;
