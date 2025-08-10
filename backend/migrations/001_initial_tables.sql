-- +goose Up
create table "user"
(
    id             serial primary key,
    name           varchar not null
);

create table trip
(
    id             serial primary key,
    name           varchar(255) not null,
    start_date     date not null,
    end_date       date not null,
    description    varchar
);

create table location
(
    id             serial primary key,
    latitude       real not null,
    longitude      real not null
);

create table activity
(
    id             serial primary key,
    trip_id        integer not null references trip on delete cascade,
    location_id    integer references location on delete set null,
    name           varchar(255) not null,
    date           date not null,
    time           time,
    description    varchar,
    address        varchar(255),
    price          integer
);

create table accommodation
(
    id             serial primary key,
    trip_id        integer not null references trip on delete cascade,
    location_id    integer references location on delete set null,
    name           varchar(255) not null,
    arrival_date   date not null,
    departure_date date not null,
    check_in_time  time,
    check_out_time time,
    description    varchar,
    address        varchar(255),
    price          integer
);

create table flight
(
    id             serial primary key,
    trip_id        integer not null references trip on delete cascade,
    price          integer
--     type           varchar(20) not null,
--     geojson        json,
--     price          integer
);

create table airport
(
    iata           varchar(3) primary key,
    name           varchar(255) not null,
    municipality   varchar(255) not null,
    location_id    integer references location on delete set null
);

create table flight_leg
(
    id                  serial primary key,
    flight_id           integer not null references flight on delete cascade,
    origin              varchar(3) not null references airport,
    destination         varchar(3) not null references airport,
    airline             varchar(255) not null,
    flight_number       varchar(255) not null,
    departure_time      timestamp not null,
    arrival_time        timestamp not null,
    duration_in_minutes integer not null,
    aircraft            varchar(255)
);

create table pnr
(
    id             serial primary key,
    flight_id      integer not null references flight on delete cascade,
    airline        varchar(10) not null,
    pnr            varchar(10) not null
);

-- +goose Down
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS trip;
DROP TABLE IF EXISTS activity;
DROP TABLE IF EXISTS accomodation;
DROP TABLE IF EXISTS flight;
DROP TABLE IF EXISTS airport;
DROP TABLE IF EXISTS flight_leg;
DROP TABLE IF EXISTS pnr;
