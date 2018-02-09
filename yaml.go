package kit

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ReadYAMLFile reads a yaml file from disk
func ReadYAMLFile(filename string, v interface{}) error {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(data, v); err != nil {
		return err
	}
	return err
}
