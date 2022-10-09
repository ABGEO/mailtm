package util

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransformerSuite struct {
	suite.Suite
}

func TestTransformerSuiteSuite(t *testing.T) {
	suite.Run(t, new(TransformerSuite))
}

func (suite *HTTPErrorSuite) TestEmailAddressesToString_OneWithoutName() {
	address1 := dto.EmailAddress{
		Address: "foo@bar.baz",
	}
	transformed := EmailAddressesToString(address1)

	assert.Equal(suite.T(), "foo@bar.baz", transformed)
}

func (suite *HTTPErrorSuite) TestEmailAddressesToString_OneWithName() {
	address1 := dto.EmailAddress{
		Address: "foo@bar.baz",
		Name:    "Foo Bar",
	}
	transformed := EmailAddressesToString(address1)

	assert.Equal(suite.T(), "foo@bar.baz (Foo Bar)", transformed)
}

func (suite *HTTPErrorSuite) TestEmailAddressesToString_ManyWithoutName() {
	address1 := dto.EmailAddress{
		Address: "foo@bar.baz",
	}
	address2 := dto.EmailAddress{
		Address: "baz@bar.foo",
	}
	transformed := EmailAddressesToString(address1, address2)

	assert.Equal(suite.T(), "foo@bar.baz, baz@bar.foo", transformed)
}

func (suite *HTTPErrorSuite) TestEmailAddressesToString_ManyWithName() {
	address1 := dto.EmailAddress{
		Address: "foo@bar.baz",
		Name:    "Foo Bar",
	}
	address2 := dto.EmailAddress{
		Address: "baz@bar.foo",
		Name:    "Baz Bar",
	}
	transformed := EmailAddressesToString(address1, address2)

	assert.Equal(suite.T(), "foo@bar.baz (Foo Bar), baz@bar.foo (Baz Bar)", transformed)
}

func (suite *HTTPErrorSuite) TestEmailAddressesToString_ManyMixed() {
	address1 := dto.EmailAddress{
		Address: "foo@bar.baz",
	}
	address2 := dto.EmailAddress{
		Address: "baz@bar.foo",
		Name:    "Baz Bar",
	}
	transformed := EmailAddressesToString(address1, address2)

	assert.Equal(suite.T(), "foo@bar.baz, baz@bar.foo (Baz Bar)", transformed)
}
