package config

type Config struct {
	DBUsername string `envconfig:"DB_USER" default:"root"`
	DBPassword string `envconfig:"DB_PASS" default:"Pastibisa"`
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconfig:"DB_PORT" default:"3302"`
	DBName     string `envconfig:"DB_NAME" default:"Gin_todo"`
}
