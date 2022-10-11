package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RandomSuite struct {
	suite.Suite
}

func TestRandomSuiteSuite(t *testing.T) {
	suite.Run(t, new(RandomSuite))
}

func (suite *RandomSuite) TestRandomString() {
	random := RandomString(16)

	assert.Equal(suite.T(), 16, len(random))
	assert.Regexp(suite.T(), `^[a-zA-Z0-9]*$`, random)
}
