package errors

import (
	"fmt"
	"os"
)

var errHandler = func(msg interface{}, code int) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(code)
}

func SetErrorHandler(handler func(msg interface{}, code int)) {
	errHandler = handler
}

func CheckErr(msg interface{}, code int) {
	if msg != nil {
		errHandler(msg, code)
	}
}
