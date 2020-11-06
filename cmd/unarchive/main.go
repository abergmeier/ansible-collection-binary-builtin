package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/abergmeier/ansible-collection-binary-builtin/internal/protocol"
	"github.com/abergmeier/ansible-collection-binary-builtin/internal/unarchive"
)

func main() {

	if len(os.Args) != 2 {
		protocol.PrintFailed("No argument file provided")
		os.Exit(1)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		protocol.PrintfFailed("Could not read configuration file %s: %s", argsFile, err)
		os.Exit(2)
	}

	p := &unarchive.Parameters{}
	err = json.Unmarshal(text, p)
	if err != nil {
		protocol.PrintfFailed("Configuration file not valid JSON %s: %s", argsFile, err)
		os.Exit(4)
	}

	err = unarchive.Run(p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(8)
	}

	protocol.PrintfSuccess("Checked out successfully")
}
