package protocol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ReadParametersFromFile reads the Ansible parameter file and unmarshalls
func ReadParametersFromFile(path string, v interface{}) error {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Could not read configuration file %s: %s", path, err)
	}

	err = json.Unmarshal(text, v)
	if err != nil {
		return fmt.Errorf("Configuration file not valid JSON %s: %s", path, err)
	}
	return nil
}
