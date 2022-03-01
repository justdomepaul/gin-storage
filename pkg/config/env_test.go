package config

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type EnvSetSuite struct {
	suite.Suite
}

func (suite *EnvSetSuite) TestLoadFromEnv() {
	core := Core{}
	suite.NoError(LoadFromEnv(&core))
}

func TestEnvSetSuite(t *testing.T) {
	suite.Run(t, new(EnvSetSuite))
}
