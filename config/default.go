package config

func defaultConfig() *Config {
	defaultCfg := &Config{}

	// Server config
	defaultCfg.Server.Port = 8080
	defaultCfg.Server.Host = "localhost"

	// Database config
	defaultCfg.Database.File = "database.db"

	// Task managers
	defaultCfg.TaskManagers = []ManagerConfig{
		{
			Name:          "default",
			DisplayName:   "Default",
			ActivePath:    "./tasks/default/active",
			CompletedPath: "./tasks/default/completed",
		},
	}

	return defaultCfg
}
