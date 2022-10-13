package errors

import (
	"fmt"
	"os"
)

var errHandler = GetDefaultErrorHandler()

func GetDefaultErrorHandler() func(msg interface{}, code int) {
	return func(msg interface{}, code int) {
		fmt.Fprintln(os.Stderr, "Error:", msg)
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
