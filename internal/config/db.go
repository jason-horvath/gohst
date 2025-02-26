package config

type DatabaseConfig struct {
	Host	 	string
	Port	 	int
	User	string
	Password 	string
	DBName		string
}

var DB *DatabaseConfig

const DB_DEFAULT_PORT = 3306

func initDB() {
	DB = &DatabaseConfig{
		Host: GetEnv("DB_HOST", "localhost").(string),
		Port: GetEnv("DB_PORT", DB_DEFAULT_PORT).(int),
		User: GetEnv("DB_USER", "default").(string),
		Password: GetEnv("DB_PASSWORD", "password").(string),
		DBName: GetEnv("DB_NAME", "db_name").(string),
	}
}

