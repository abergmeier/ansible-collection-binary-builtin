package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/abergmeier/ansible-collection-binary-builtin/internal/git"
)

type response struct {
	Msg     string `json:"msg"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

func main() {

	if len(os.Args) != 2 {
		printFailed("No argument file provided")
		os.Exit(1)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		printfFailed("Could not read configuration file %s: %s", argsFile, err)
		os.Exit(2)
	}

	p := &git.Parameters{}
	err = json.Unmarshal(text, p)
	if err != nil {
		printfFailed("Configuration file not valid JSON %s: %s", argsFile, err)
		os.Exit(4)
	}

	err = git.Run(p)
	if err != nil {
		printFailed(err.Error())
		os.Exit(8)
	}

	printfSuccess("Checked out successfully")
}

func printSuccess(msg string) {
	jr, err := json.Marshal(response{
		Msg: msg,
	})
	if err != nil {
		jr, _ = json.Marshal(response{Msg: "Invalid response object"})
	}
	fmt.Println(string(jr))
}

func printfSuccess(format string, a ...interface{}) {
	printSuccess(fmt.Sprintf(format, a...))
}

func printFailed(msg string) {

	jr, err := json.Marshal(response{
		Msg:    msg,
		Failed: true,
	})
	if err != nil {
		jr, _ = json.Marshal(response{Msg: "Invalid response object"})
	}
	fmt.Println(string(jr))
}

func printfFailed(format string, a ...interface{}) {
	printFailed(fmt.Sprintf(format, a...))
}
