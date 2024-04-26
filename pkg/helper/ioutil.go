package helper

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadFile2Obj(filename string, obj interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read config(file:%s): %w", filename, err)
	}

	err = yaml.Unmarshal(file, obj)
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}
