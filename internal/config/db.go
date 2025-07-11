package config

const DB_DEFAULT_PORT = 5432
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

// Package config provides configuration management for the application, including database configurations.
func NewDatabaseConfigPool() *DatabaseConfigPool {
	return &DatabaseConfigPool{
		configs: make(map[string]*DatabaseConfig),
	}
}

// NewDatabaseConfigPool creates a new instance of DatabaseConfigPool.
func (p *DatabaseConfigPool) Add(name string, config *DatabaseConfig) {
	p.configs[name] = config
}

// AddOrUpdate adds a new database configuration or updates an existing one.
func (p *DatabaseConfigPool) Get(name string) (*DatabaseConfig, bool) {
	config, exists := p.configs[name]
	return config, exists
}

// GetConfigs returns all database configurations in the pool.
func (p *DatabaseConfigPool) GetConfigs() map[string]*DatabaseConfig {
	return p.configs
}

// GetOrDefault retrieves a database configuration by name, returning a default configuration if it doesn't exist.
func (p *DatabaseConfigPool) GetOrDefault(name string, defaultConfig *DatabaseConfig) *DatabaseConfig {
	if config, exists := p.configs[name]; exists {
		return config
	}
	return defaultConfig
}

// Remove deletes a database configuration by name.
func (p *DatabaseConfigPool) Has(name string) bool {
	_, exists := p.configs[name]
	return exists
}

// Remove deletes a database configuration by name.
func (p *DatabaseConfigPool) Names() []string {
	names := make([]string, 0, len(p.configs))
	for name := range p.configs {
		names = append(names, name)
	}
	return names
}
