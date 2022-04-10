package config

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func Init() (*PostgresConfig, error) {
	return nil, nil
}

func parseConfigFile() error {
	return nil
}

func unmarshal(cfg *PostgresConfig) error {
	return nil
}
