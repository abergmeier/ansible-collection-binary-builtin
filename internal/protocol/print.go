package protocol

import (
	"encoding/json"
	"fmt"
)

// PrintSuccess prints a message to stdout in JSON format
func PrintSuccess(msg string) {
	jr, err := json.Marshal(Response{
		Msg: msg,
	})
	if err != nil {
		jr, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(jr))
}

// PrintfSuccess prints a formatted message to stdout in JSON format
func PrintfSuccess(format string, a ...interface{}) {
	PrintSuccess(fmt.Sprintf(format, a...))
}

// PrintFailed prints a fail message to stdout in JSON format
func PrintFailed(msg string) {

	jr, err := json.Marshal(Response{
		Msg:    msg,
		Failed: true,
	})
	if err != nil {
		jr, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(jr))
}

// PrintfFailed prints a formatted fail message to stdout in JSON format
func PrintfFailed(format string, a ...interface{}) {
	PrintFailed(fmt.Sprintf(format, a...))
}
