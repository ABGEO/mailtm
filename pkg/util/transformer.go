package util

import (
	"fmt"
	"strings"

	"github.com/abgeo/mailtm/pkg/dto"
)

func EmailAddressesToString(addresses ...dto.EmailAddress) string {
	formattedAddresses := make([]string, len(addresses))

	for i, address := range addresses {
		formattedAddress := address.Address
		if address.Name != "" {
			formattedAddress += fmt.Sprintf("(%s)", address.Name)
		}

		formattedAddresses[i] = formattedAddress
	}

	return strings.Join(formattedAddresses, ", ")
}
