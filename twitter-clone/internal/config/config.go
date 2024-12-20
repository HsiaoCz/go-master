package config

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	Server struct {
		Port int
	}
}

func LoadConfig() *Config {
	// In a real application, you'd load this from environment variables
	// or a configuration file
	cfg := &Config{}
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.User = "postgres"
	cfg.Database.Password = "postgres"
	cfg.Database.Name = "twitter_clone"
	cfg.Server.Port = 8080
	return cfg
}