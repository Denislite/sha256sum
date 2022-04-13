package hash

type FileInfo struct {
	FileName  string `db:"file_name"`
	FilePath  string `db:"file_path"`
	HashValue []byte `db:"hash_value"`
	HashType  string `db:"hash_type"`
}
