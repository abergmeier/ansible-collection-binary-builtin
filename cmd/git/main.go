package main

import (
	"os"

	"github.com/abergmeier/ansible-collection-binary-builtin/internal/git"
	"github.com/abergmeier/ansible-collection-binary-builtin/internal/protocol"
)

func main() {

	if len(os.Args) != 2 {
		protocol.PrintFailed("No argument file provided")
		os.Exit(1)
	}

	argsFile := os.Args[1]
	p := &git.Parameters{}
	err := protocol.ReadParametersFromFile(argsFile, p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(2)
	}

	err = git.Run(p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(4)
	}

	protocol.PrintfSuccess("Checked out successfully")
}
