package messager

import (
	"fmt"
	"io"
	"os"
	"encoding/json"
	"github.com/mitchellh/colorstring"
)

type ResourceMessager struct {
	logWriter      io.Writer
	responseWriter io.Writer
	exitOnFatal    bool
}

var logger *ResourceMessager

func NewMessager(logWriter, responseWriter io.Writer) (*ResourceMessager) {
	return &ResourceMessager{logWriter, responseWriter, true}
}
func (rl *ResourceMessager) LogIt(args ...interface{}) {
	var text string
	text, ok := args[0].(string);
	if !ok {
		panic("Firt argument should be a string")
	}
	text = colorstring.Color(text)
	if len(args) > 1 {
		newArgs := args[1:]
		fmt.Fprintf(rl.logWriter, text, newArgs...)
	} else {
		fmt.Fprint(rl.logWriter, text)
	}
}
func (rl *ResourceMessager) LogItLn(args ...interface{}) {
	var text string
	text, ok := args[0].(string);
	if !ok {
		panic("Firt argument should be a string")
	}
	args[0] = text + "\n"
	rl.LogIt(args...)
}
func (rl *ResourceMessager) SendJsonResponse(v interface{}) {
	json.NewEncoder(rl.responseWriter).Encode(v)
}
func (rl *ResourceMessager) GetLogWriter() (io.Writer) {
	return rl.logWriter
}
func (rl *ResourceMessager) GetResponseWriter() (io.Writer) {
	return rl.responseWriter
}
func (rl *ResourceMessager) FatalIf(doing string, err error) {
	if err != nil {
		rl.Fatal(doing + ": " + err.Error())
	}
}
func (rl *ResourceMessager) Fatal(message string) {
	fmt.Fprintln(rl.responseWriter, message)
	if rl.exitOnFatal {
		os.Exit(1)
	}
}
func (rl *ResourceMessager) SetExitOnFatal(exitOnFatal bool) {
	rl.exitOnFatal = exitOnFatal
}
func GetMessager() (*ResourceMessager) {
	if logger == nil {
		logger = NewMessager(os.Stderr, os.Stdout)
	}
	return logger
}