CREATE TABLE files
(
    id          BIGSERIAL PRIMARY KEY,
    file_name   VARCHAR,
    file_path   TEXT,
    hash_value  VARCHAR,
    hash_type   VARCHAR,
    deleted     BOOLEAN DEFAULT FALSE,

    CONSTRAINT files_unique UNIQUE (file_path, hash_type)
);

INSERT INTO files (file_name, file_path, hash_value, hash_type) VALUES
    ('1.txt','1/1.txt','123','sha256'),
    ('1.txt','1/1.txt','1234','sha512'),
    ('2.txt','1/2.txt','123','sha256'),
    ('3.txt','1/3.txt','123','sha256'),
    ('3.txt','1/3.txt','1234','sha512');