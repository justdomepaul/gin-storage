package config

// Media type
type Media struct {
	StorageDomain string `split_words:"true" default:"http://localhost"`
	BucketName    string `split_words:"true" default:"staging.megaphone.appspot.com"`
	PrefixPath    string `split_words:"true" default:""`
}
