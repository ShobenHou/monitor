package yaml

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadYaml(path string, out interface{}) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, out)
	if err != nil {
		return err
	}
	return nil
}
