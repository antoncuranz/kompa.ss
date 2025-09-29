-- +goose Up
ALTER TABLE train_detail
ALTER COLUMN refresh_token TYPE varchar;