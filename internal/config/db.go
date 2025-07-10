package config

type DatabaseConfig struct {
	Host	 	string
	Port	 	int
	User	string
	Password 	string
	DBName		string
}

type DatabaseConfigPool struct {
	configs map[string]*DatabaseConfig
}

func NewDatabaseConfigPool() *DatabaseConfigPool {
	return &DatabaseConfigPool{
		configs: make(map[string]*DatabaseConfig),
	}
}

func (p *DatabaseConfigPool) Add(name string, config *DatabaseConfig) {
	p.configs[name] = config
}

func (p *DatabaseConfigPool) Get(name string) (*DatabaseConfig, bool) {
	config, exists := p.configs[name]
	return config, exists
}

func (p *DatabaseConfigPool) GetOrDefault(name string, defaultConfig *DatabaseConfig) *DatabaseConfig {
	if config, exists := p.configs[name]; exists {
		return config
	}
	return defaultConfig
}

func (p *DatabaseConfigPool) Has(name string) bool {
	_, exists := p.configs[name]
	return exists
}

func (p *DatabaseConfigPool) Names() []string {
	names := make([]string, 0, len(p.configs))
	for name := range p.configs {
		names = append(names, name)
	}
	return names
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

