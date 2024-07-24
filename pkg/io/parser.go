package io

import (
	"fmt"
	"io/ioutil"
)

type Parser interface {
	Load(string, interface{}) error
	Save(string) error
}
type parser struct {
	data interface{}
}

func ReadFile(file string) (data []byte, err error) {
	if len(file) < 1 {
		return nil, fmt.Errorf("invalid file: %s\n%s", file, err)
	}

	// Read in the files
	data, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s\n%s", file, err)
	}

	return data, nil
}
