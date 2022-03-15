package config

import (
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type MediaSuite struct {
	suite.Suite
	StorageDomain string
	BucketName    string
}

func (suite *MediaSuite) SetupSuite() {
	t := suite.T()
	os.Clearenv()
	suite.StorageDomain = "http://testDomain.com"
	suite.BucketName = "testBucket"
	assert.NoError(t, os.Setenv("STORAGE_DOMAIN", suite.StorageDomain))
	assert.NoError(t, os.Setenv("BUCKET_NAME", suite.BucketName))
}

func (suite *MediaSuite) TestDefaultOption() {
	t := suite.T()
	options := &Media{}
	suite.NoError(config.LoadFromEnv(options))
	assert.Equal(t, suite.StorageDomain, options.StorageDomain)
	assert.Equal(t, suite.BucketName, options.BucketName)
}

func TestMediaSuite(t *testing.T) {
	suite.Run(t, new(MediaSuite))
}
