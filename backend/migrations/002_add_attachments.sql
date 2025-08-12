-- +goose Up
create table attachment
(
    id               serial primary key,
    trip_id          integer not null references trip on delete cascade,
    name             varchar(255) not null,
    blob             bytea not null
);

create table attachment_activity
(
    attachment_id    integer not null references attachment on delete cascade,
    activity_id      integer not null references activity on delete cascade
);

create table attachment_accommodation
(
    attachment_id    integer not null references attachment on delete cascade,
    accommodation_id integer not null references accommodation on delete cascade
);

create table attachment_flight
(
    attachment_id    integer not null references attachment on delete cascade,
    flight_id        integer not null references flight on delete cascade
);

-- +goose Down
DROP TABLE IF EXISTS attachment;
