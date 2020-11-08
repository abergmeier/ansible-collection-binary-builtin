package main

import (
	"os"

	"github.com/abergmeier/ansible-collection-binary-builtin/internal/assert"
	"github.com/abergmeier/ansible-collection-binary-builtin/internal/protocol"
)

func main() {

	if len(os.Args) != 2 {
		protocol.PrintFailed("No argument file provided")
		os.Exit(1)
	}

	argsFile := os.Args[1]
	p := &assert.Parameters{}
	err := protocol.ReadParametersFromFile(argsFile, p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(2)
	}

	truth, err := assert.Run(p)
	if err != nil {
		protocol.PrintFailed(err.Error())
		os.Exit(4)
	}

	if !truth {
		protocol.PrintfSuccess(p.FailMsg)
		os.Exit(8)
	}

	protocol.PrintfSuccess(p.SuccessMsg)
}
