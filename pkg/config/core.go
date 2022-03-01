package config

// Core type
type Core struct {
	LoggerMode  string `split_words:"true" default:"customized"`
	ReleaseMode bool   `split_words:"true" default:"true"`
	SystemName  string `split_words:"true" default:"system"`
}
