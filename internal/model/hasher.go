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

type DeletedFiles struct {
	FileName string
	OldHash  string
	FilePath string
}
