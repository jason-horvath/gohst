package config


type FileConfig struct {
	Path 	string
}
type RedisConfig struct {
	DB			int
	Host 		string
	Password 	string
	Port 		int
}

type SessionConfig struct {
	ContextKey 	string
	File 		*FileConfig
	Length 		int
	Name 		string
	Redis 		*RedisConfig
	Store 		string
}

const SESSION_LENGTH_DEFAULT = 60

var Session *SessionConfig

func initSession() {
	Session = &SessionConfig{
		ContextKey: GetEnv("SESSION_CONTEXT_KEY" , "session").(string),
		File: &FileConfig{
			Path: GetEnv("SESSION_FILE_PATH", "tmp/sessions").(string),
		},
		Length: GetEnv("SESSION_LENGTH", SESSION_LENGTH_DEFAULT).(int),
		Name: GetEnv("SESSION_NAME", "session_id").(string),
		Redis: &RedisConfig{
			DB: GetEnv("SESSION_REDIS_DB", 0).(int),
			Host: GetEnv("SESSION_REDIS_HOST", "localhost").(string),
			Password: GetEnv("SESSION_REDIS_PASSWORD", "").(string),
			Port: GetEnv("SESSION_REDIS_PORT", 6379).(int),
		},
		Store: GetEnv("SESSION_STORE", "file").(string),
	}
}
