-- +goose Up
create table transportation_generic
(
    transportation_id   integer primary key references transportation on delete cascade,
    name                varchar(255) not null,
    origin_address      varchar(255),
    destination_address varchar(255)
);

-- +goose Down
DROP TABLE IF EXISTS transportation_generic;
