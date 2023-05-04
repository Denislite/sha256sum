CREATE TABLE files
(
    id              UUID PRIMARY KEY,
    pod_name        VARCHAR,
    container_id    UUID,
    image_name      VARCHAR,
    image_version   VARCHAR,
    file_name       VARCHAR,
    file_path       TEXT,
    hash_value      VARCHAR,
    hash_type       VARCHAR,
    created_at      TIMESTAMP DEFAULT now(),

    CONSTRAINT files_unique UNIQUE (file_path, hash_type),
    FOREIGN KEY (container_id) REFERENCES containers (id)
);

CREATE TABLE containers
(
    id UUID PRIMARY KEY,
    container_name VARCHAR,
    image_name VARCHAR,
    image_version VARCHAR
);

CREATE TABLE verifications
(
    id UUID PRIMARY KEY,
    info_id UUID,
    container_name VARCHAR,
    image_name VARCHAR,
    image_version VARCHAR,
    FOREIGN KEY (info_id) REFERENCES verificationInfo (id)
);

CREATE TABLE verificationInfo
(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    status VARCHAR
);

CREATE TABLE verificationFiles
(
    id_verification UUID PRIMARY KEY,
    id_file UUID PRIMARY KEY,
    FOREIGN KEY (id_file) REFERENCES fileInfo (id),
    FOREIGN KEY (id_verification) REFERENCES verificationInfo (id)
);

CREATE TABLE fileInfo
(
    id UUID PRIMARY KEY,
    file_name       VARCHAR,
    hash_value      VARCHAR,
);