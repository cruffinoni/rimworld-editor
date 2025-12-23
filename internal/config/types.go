package config

type Config struct {
	Logging LoggingConfig `mapstructure:"logging"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level" validate:"oneof=trace debug info warn error fatal"`
	Format string `mapstructure:"format" validate:"oneof=json text"`
	Output string `mapstructure:"output"`
}
