package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type CloudSuite struct {
	suite.Suite
	ProjectID             string
	EndPoint              string
	WithoutAuthentication bool
	GRPCInsecure          bool
}

func (suite *CloudSuite) SetupSuite() {
	t := suite.T()
	os.Clearenv()
	suite.ProjectID = "testProjectID"
	suite.EndPoint = "testEndPoint"
	suite.WithoutAuthentication = true
	suite.GRPCInsecure = true

	assert.NoError(t, os.Setenv("PROJECT_ID", suite.ProjectID))
	assert.NoError(t, os.Setenv("END_POINT", suite.EndPoint))
	assert.NoError(t, os.Setenv("WITHOUT_AUTHENTICATION", strconv.FormatBool(suite.WithoutAuthentication)))
	assert.NoError(t, os.Setenv("GRPC_INSECURE", strconv.FormatBool(suite.GRPCInsecure)))

}

func (suite *CloudSuite) TestDefaultOption() {
	t := suite.T()
	spanner := &Cloud{}
	suite.NoError(LoadFromEnv(spanner))
	assert.Equal(t, suite.ProjectID, spanner.ProjectID)
	assert.Equal(t, suite.EndPoint, spanner.EndPoint)
	assert.Equal(t, suite.WithoutAuthentication, spanner.WithoutAuthentication)
	assert.Equal(t, suite.GRPCInsecure, spanner.GRPCInsecure)
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(CloudSuite))
}
