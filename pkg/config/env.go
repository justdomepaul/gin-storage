package config

import "github.com/kelseyhightower/envconfig"

func LoadFromEnv(pointer ...interface{}) error {
	for _, p := range pointer {
		err := envconfig.Process("", p)
		if err != nil {
			return err
		}
	}
	return nil
}
