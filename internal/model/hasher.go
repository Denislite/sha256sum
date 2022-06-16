package model

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ChangedFiles struct {
	FileName string
	OldHash  string
	NewHash  string
}

type FileInfo struct {
	FileName  string `db:"file_name"`
	FilePath  string `db:"file_path"`
	HashValue string `db:"hash_value"`
	HashType  string `db:"hash_type"`
	Deleted   bool   `db:"deleted"`
}
