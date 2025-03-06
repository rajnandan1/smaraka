package models

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type StartupConfig struct {
	Port                    string `yaml:"OK_BOOKMARK_port"`
	Environment             string `yaml:"OK_BOOKMARK_ENV"`
	MaxWorkers              int    `yaml:"OK_BOOKMARK_MaxWorkers"`
	GracefulShutDownTimeout int    `yaml:"OK_BOOKMARK_GracefulShutDownTimeout"`

	StytchProjectID       string `yaml:"OK_BOOKMARK_StytchProjectID"`
	StytchSecret          string `yaml:"OK_BOOKMARK_StytchSecret"`
	StytchSessionDuration int    `yaml:"OK_BOOKMARK_StytchSessionDuration"`
	StytchPublicToken     string `yaml:"OK_BOOKMARK_StytchPublicToken"`

	PostgresUser     string `yaml:"OK_BOOKMARK_PostgresUser"`
	PostgresPassword string `yaml:"OK_BOOKMARK_PostgresPassword"`
	PostgresHost     string `yaml:"OK_BOOKMARK_PostgresHost"`
	PostgresPort     string `yaml:"OK_BOOKMARK_PostgresPort"`
	PostgresDB       string `yaml:"OK_BOOKMARK_PostgresDB"`
}
