CREATE TABLE files
(
    id          BIGSERIAL PRIMARY KEY,
    file_name   VARCHAR,
    file_path   TEXT,
    hash_value  VARCHAR,
    hash_type   VARCHAR
);

CREATE FUNCTION check_hash (VARCHAR, TEXT, VARCHAR, VARCHAR) RETURNS void
AS $$
BEGIN
   IF NOT EXISTS (SELECT * FROM files
				  WHERE file_path=$2
                  AND hash_type=$4)
   THEN
        INSERT INTO files (file_name,file_path,hash_value,hash_type)
        VALUES ($1, $2, $3, $4);
   ELSE
   		UPDATE files SET hash_value=$3 WHERE file_path=$2 AND hash_type=$4;
   END IF;
END;
$$ LANGUAGE plpgsql;