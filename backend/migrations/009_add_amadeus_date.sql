-- +goose Up
alter table flight_leg add column amadeus_date date;

-- +goose Down
alter table flight_leg drop column amadeus_date;
