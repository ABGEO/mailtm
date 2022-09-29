package util

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type HTTPError struct {
	Code       int                 `json:"code,omitempty"`
	Type       string              `json:"type"`
	Title      string              `json:"title"`
	Detail     string              `json:"detail"`
	Violations []map[string]string `json:"violations"`
}

func NewHTTPError(response *resty.Response) (err *HTTPError) {
	err = response.Error().(*HTTPError)
	err.Code = response.StatusCode()

	if err.Detail == "" {
		pErr := struct {
			Message string `json:"message"`
		}{}

		if mErr := json.Unmarshal(response.Body(), &pErr); mErr == nil && pErr.Message != "" {
			err.Type = "https://www.rfc-editor.org/rfc/rfc2616#section-10"
			err.Title = "An error occurred"
			err.Detail = pErr.Message
		}
	}

	return err
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("%s [%d]", err.Detail, err.Code)
}
