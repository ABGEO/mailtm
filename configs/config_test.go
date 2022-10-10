package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestNewConfig_CreateIfFileNotFound" {
		if homedir, err := os.UserHomeDir(); err == nil {
			_ = os.RemoveAll(homedir + "/.mailtm")
		}
	}
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (suite *ConfigSuite) TestNewConfig_CreateIfFileNotFound() {
	_ = NewConfig()

	assert.FileExists(suite.T(), suite.getConfigDir()+"/config")
}

func (suite *ConfigSuite) TestWrite() {
	oldConfig := NewConfig()
	oldConfig.Auth.ID = "foo"
	oldConfig.Auth.Email = "bar"
	oldConfig.Auth.Token = "baz"

	oldConfig.Write()

	newConfig := NewConfig()
	assert.Equal(suite.T(), "foo", newConfig.Auth.ID)
	assert.Equal(suite.T(), "bar", newConfig.Auth.Email)
	assert.Equal(suite.T(), "baz", newConfig.Auth.Token)
}

func (suite *ConfigSuite) getConfigDir() string {
	if homedir, err := os.UserHomeDir(); err == nil {
		return homedir + "/.mailtm"
	}

	return ""
}
