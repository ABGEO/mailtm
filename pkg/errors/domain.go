package errors

import (
	"fmt"
	"strings"
)

type InvalidDomainError struct {
	Domain       string
	ValidDomains []string
}

func NewInvalidDomainError(domains []string, validDomain string) *InvalidDomainError {
	return &InvalidDomainError{
		Domain:       validDomain,
		ValidDomains: domains,
	}
}

func (err *InvalidDomainError) Error() string {
	return fmt.Sprintf(
		"domain %s is not valid. Valid domains are: [%s]",
		err.Domain,
		strings.Join(err.ValidDomains, ", "))
}
