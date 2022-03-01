package config

// Cloud type
type Cloud struct {
	ProjectID             string `split_words:"true" default:"test-project"`
	EndPoint              string `split_words:"true" default:"localhost:9010"`
	WithoutAuthentication bool   `split_words:"true" default:"true"`
	GRPCInsecure          bool   `split_words:"true" default:"true"`
}
