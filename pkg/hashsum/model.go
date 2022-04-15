package hashsum

type FileInfo struct {
	FileName  string `db:"file_name"`
	FilePath  string `db:"file_path"`
	HashValue string `db:"hash_value"`
	HashType  string `db:"hash_type"`
	Deleted   bool   `db:"deleted"`
}
