package storage

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type testIFile struct {
	mock.Mock
	IFile
}

type DriverSuite struct {
	suite.Suite
}

func (suite *DriverSuite) TestRegister() {
	testIFile := &testIFile{}
	Register(testIFile, func() {})
	defer Unload()
	suite.Equal("*storage.testIFile", reflect.TypeOf(FILE).String())
	suite.Equal("func()", reflect.TypeOf(CLOSEFN).String())
}

func (suite *DriverSuite) TestLoad() {
	testIFile := &testIFile{}
	Register(testIFile, func() {})
	defer Unload()
	f, cn := Load()
	suite.Equal("*storage.testIFile", reflect.TypeOf(f).String())
	suite.Equal("func()", reflect.TypeOf(cn).String())
}

func (suite *DriverSuite) TestLoadError() {
	testIFile := &testIFile{}
	Register(testIFile, func() {})
	Unload()
	suite.Panics(func() {
		Load()
	})
}

func (suite *DriverSuite) TestUnload() {
	Unload()
	suite.Nil(FILE)
	suite.Nil(CLOSEFN)
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(DriverSuite))
}
