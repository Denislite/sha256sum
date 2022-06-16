CREATE TABLE files
(
    id              BIGSERIAL PRIMARY KEY,
    image_name      VARCHAR,
    image_version   VARCHAR,
    file_name       VARCHAR,
    file_path       TEXT,
    hash_value      VARCHAR,
    hash_type       VARCHAR,
    created_at      TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),

    CONSTRAINT files_unique UNIQUE (file_path, hash_type)
);

INSERT INTO files (file_name, file_path, hash_value, hash_type) VALUES
    ('1.txt','1/1.txt','123','sha256'),
    ('1.txt','1/1.txt','1234','sha512'),
    ('2.txt','1/2.txt','123','sha256'),
    ('3.txt','1/3.txt','123','sha256'),
    ('3.txt','1/3.txt','1234','sha512');