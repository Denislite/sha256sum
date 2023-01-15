CREATE TABLE files
(
    id              BIGSERIAL PRIMARY KEY,
    pod_name        VARCHAR,
    image_name      VARCHAR,
    image_version   VARCHAR,
    file_name       VARCHAR,
    file_path       TEXT,
    hash_value      VARCHAR,
    hash_type       VARCHAR,
    created_at      TIMESTAMP DEFAULT now(),

    CONSTRAINT files_unique UNIQUE (file_path, hash_type)
);