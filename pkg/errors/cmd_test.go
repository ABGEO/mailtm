package errors

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CmdErrorSuite struct {
	suite.Suite
}

func (suite *CmdErrorSuite) TearDownTest() {
	SetDefaultErrorHandler()
}

func TestCmdErrorSuite(t *testing.T) {
	suite.Run(t, new(CmdErrorSuite))
}

func (suite *CmdErrorSuite) TestSetErrorHandler() {
	customErrorHandler := func(msg interface{}, code int) {
		suite.T().Error(msg, code)
	}

	SetErrorHandler(customErrorHandler)

	customErrorHandlerFuncName := runtime.FuncForPC(reflect.ValueOf(customErrorHandler).Pointer()).Name()
	errHandlerFuncName := runtime.FuncForPC(reflect.ValueOf(errHandler).Pointer()).Name()

	assert.Equal(suite.T(), customErrorHandlerFuncName, errHandlerFuncName)
}

func (suite *CmdErrorSuite) TestCheckErr_WithDefaultHandler() {
	fakeExit := func(code int) { panic(fmt.Sprintf("exited with %d", code)) }
	patch := monkey.Patch(os.Exit, fakeExit)

	defer patch.Unpatch()

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { CheckErr("error", 1) })
}

func (suite *CmdErrorSuite) TestCheckErr_WithCustomHandler() {
	var errorMessage string

	SetErrorHandler(func(msg interface{}, code int) {
		errorMessage = fmt.Sprintf("%s (%d)", msg, code)
	})
	CheckErr("error", 1)

	assert.Equal(suite.T(), "error (1)", errorMessage)
}

func (suite *CmdErrorSuite) TestCheckErr_WithAccessDeniedError() {
	fakeExit := func(code int) { panic(fmt.Sprintf("exited with %d", code)) }
	patch := monkey.Patch(os.Exit, fakeExit)

	defer patch.Unpatch()

	accessDeniedError := &HTTPError{
		Code:   401,
		Detail: "Access Denied",
	}

	testStderr := bytes.NewBufferString("")
	originalStderr := os.Stderr
	readStream, writeStream, _ := os.Pipe()
	os.Stderr = writeStream
	outChain := make(chan string)

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { CheckErr(accessDeniedError, 1) })

	go func() {
		_, _ = io.Copy(testStderr, readStream)
		outChain <- testStderr.String()
	}()

	_ = writeStream.Close()

	os.Stderr = originalStderr
	out := <-outChain

	assert.Contains(suite.T(), out, "try to run the \"mailtm auth\" command")
}
