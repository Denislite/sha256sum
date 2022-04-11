CREATE TABLE files
(
    id          BIGSERIAL PRIMARY KEY,
    file_name   VARCHAR,
    file_path   TEXT,
    hash_value  VARCHAR,
    hash_type   VARCHAR
);