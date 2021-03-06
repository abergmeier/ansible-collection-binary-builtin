package main

import (
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
	p := &unarchive.Parameters{}
	err := protocol.ReadParametersFromFile(argsFile, p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(2)
	}

	err = unarchive.Run(p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(8)
	}

	protocol.PrintfSuccess("Unarchive successful")
}
