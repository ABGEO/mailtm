package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InvalidDomainErrorSuite struct {
	suite.Suite
}

func TestInvalidDomainErrorSuite(t *testing.T) {
	suite.Run(t, new(InvalidDomainErrorSuite))
}

func (suite *InvalidDomainErrorSuite) TestInvalidDomainError_WithSingleValidDomain() {
	err := NewInvalidDomainError([]string{"foo.bar"}, "baz.bar")

	assert.EqualError(suite.T(), err, "domain baz.bar is not valid. Valid domains are: [foo.bar]")
}

func (suite *InvalidDomainErrorSuite) TestInvalidDomainError_WithManyValidDomains() {
	err := NewInvalidDomainError([]string{"foo.bar", "bar.baz"}, "bar.foo")

	assert.EqualError(suite.T(), err, "domain bar.foo is not valid. Valid domains are: [foo.bar, bar.baz]")
}
