-- +goose Up
alter table train_leg
add column duration_in_minutes integer not null default 0;
alter table train_leg
alter column duration_in_minutes drop default;

alter table train_leg
add column operator_name varchar(255) not null default '';
alter table train_leg
alter column operator_name drop default;