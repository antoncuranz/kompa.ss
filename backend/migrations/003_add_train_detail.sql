-- +goose Up
create table train_detail
(
    transportation_id   integer primary key references transportation on delete cascade,
    refresh_token       varchar(255) not null
);

create table train_station
(
    id             varchar(255) primary key,
    name           varchar(255) not null,
    location_id    integer references location
);

create table train_leg
(
    id                  serial primary key,
    transportation_id   integer not null references transportation on delete cascade,
    origin              varchar(255) not null references train_station,
    destination         varchar(255) not null references train_station,
    departure_time      timestamp not null,
    arrival_time        timestamp not null,
	line_name           varchar(255) not null
);

-- +goose Down
DROP TABLE IF EXISTS train_leg;
DROP TABLE IF EXISTS train_station;
DROP TABLE IF EXISTS train_detail;
