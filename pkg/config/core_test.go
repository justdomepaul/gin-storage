package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type CoreSuite struct {
	suite.Suite
	LoggerMode  string
	ReleaseMode bool
	SystemName  string
}

func (suite *CoreSuite) SetupSuite() {
	t := suite.T()
	os.Clearenv()
	suite.LoggerMode = "debug"
	suite.ReleaseMode = true
	suite.SystemName = "system"
	assert.NoError(t, os.Setenv("LOGGER_MODE", suite.LoggerMode))
	assert.NoError(t, os.Setenv("RELEASE_MODE", strconv.FormatBool(suite.ReleaseMode)))
	assert.NoError(t, os.Setenv("SYSTEM_NAME", suite.SystemName))
}

func (suite *CoreSuite) TestDefaultOption() {
	t := suite.T()
	core := &Core{}
	suite.NoError(LoadFromEnv(core))
	assert.Equal(t, suite.LoggerMode, core.LoggerMode)
	assert.Equal(t, suite.ReleaseMode, core.ReleaseMode)
	assert.Equal(t, suite.SystemName, core.SystemName)
}

func TestCoreSuite(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}
