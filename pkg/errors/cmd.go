package errors

import (
	"fmt"
	"os"
)

var errHandler = GetDefaultErrorHandler()

func GetDefaultErrorHandler() func(msg interface{}, code int) {
	return func(msg interface{}, code int) {
		errorMsg := fmt.Sprintf("Error: %s", msg)
		if err, ok := msg.(*HTTPError); ok && (err.Code == 401 || err.Code == 403) && err.Detail != "Invalid credentials." {
			errorMsg += "\ntry to run the \"mailtm auth\" command"
		}

		fmt.Fprintln(os.Stderr, errorMsg)
		os.Exit(code)
	}
}

func SetDefaultErrorHandler() {
	errHandler = GetDefaultErrorHandler()
}

func SetErrorHandler(handler func(msg interface{}, code int)) {
	errHandler = handler
}

func CheckErr(msg interface{}, code int) {
	if msg != nil {
		errHandler(msg, code)
	}
}
