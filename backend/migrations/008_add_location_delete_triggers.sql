-- +goose Up

-- +goose statementbegin
CREATE OR REPLACE FUNCTION delete_location_when_activity_deleted()
RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM location WHERE id = OLD.location_id;
    RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose statementend
CREATE TRIGGER activity_delete_trigger
    AFTER DELETE ON activity
    FOR EACH ROW
EXECUTE FUNCTION delete_location_when_activity_deleted();

-- +goose statementbegin
CREATE OR REPLACE FUNCTION delete_location_when_accommodation_deleted()
RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM location WHERE id = OLD.location_id;
    RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose statementend
CREATE TRIGGER accommodation_delete_trigger
    AFTER DELETE ON accommodation
    FOR EACH ROW
EXECUTE FUNCTION delete_location_when_accommodation_deleted();

-- +goose statementbegin
CREATE OR REPLACE FUNCTION delete_location_when_airport_deleted()
RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM location WHERE id = OLD.location_id;
    RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose statementend
CREATE TRIGGER airport_delete_trigger
    AFTER DELETE ON airport
    FOR EACH ROW
EXECUTE FUNCTION delete_location_when_airport_deleted();

-- +goose statementbegin
CREATE OR REPLACE FUNCTION delete_location_when_train_station_deleted()
RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM location WHERE id = OLD.location_id;
    RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose statementend
CREATE TRIGGER train_station_delete_trigger
    AFTER DELETE ON train_station
    FOR EACH ROW
EXECUTE FUNCTION delete_location_when_train_station_deleted();

-- +goose statementbegin
CREATE OR REPLACE FUNCTION delete_location_when_transportation_deleted()
RETURNS TRIGGER AS
$$
BEGIN
    DELETE FROM location WHERE id = OLD.origin_id;
    DELETE FROM location WHERE id = OLD.destination_id;
    RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose statementend
CREATE TRIGGER transportation_delete_trigger
    AFTER DELETE ON transportation
    FOR EACH ROW
EXECUTE FUNCTION delete_location_when_transportation_deleted();

-- +goose Down
DROP TRIGGER IF EXISTS activity_delete_trigger ON activity;
DROP FUNCTION IF EXISTS delete_location_when_activity_deleted;

DROP TRIGGER IF EXISTS accommodation_delete_trigger ON activity;
DROP FUNCTION IF EXISTS delete_location_when_accommodation_deleted;

DROP TRIGGER IF EXISTS airport_delete_trigger ON airport;
DROP FUNCTION IF EXISTS delete_location_when_airport_deleted;

DROP TRIGGER IF EXISTS train_station_delete_trigger ON train_station;
DROP FUNCTION IF EXISTS delete_location_when_train_station_deleted;

DROP TRIGGER IF EXISTS transportation_delete_trigger ON transportation;
DROP FUNCTION IF EXISTS delete_location_when_transportation_deleted;
