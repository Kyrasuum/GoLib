package io

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonParser parser

func (jp JsonParser) Save(file string) (err error) {
	if len(file) < 1 {
		return fmt.Errorf("invalid json file\n%s", err)
	}

	// convert object into a json bytelike object
	jsonObj, err := jp.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal object\n%s", err)
	}

	// write json bytelike object to file
	err = ioutil.WriteFile(file, jsonObj, 0664)
	if err != nil {
		return fmt.Errorf("failed to write to json file\n%s", err)
	}

	return err
}

func (jp JsonParser) Load(file string, structobj interface{}) (err error) {
	data, err := ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error opening json file\n%s", err)
	}

	// Parse the json into structs
	jp.data = structobj
	err = jp.FromJson(data)
	if err != nil {
		return fmt.Errorf("failed to parse json file\n%s", err)
	}

	return err
}

func (jp JsonParser) ToMap() (map[string]interface{}, error) {
	temp, err := jp.ToJson()
	if err != nil {
		return nil, err
	}
	return jp.data.(map[string]interface{}), jp.FromJson(temp)
}

func (jp JsonParser) ToJson() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(jp.data)

	return buffer.Bytes(), err
}

func (jp JsonParser) FromJson(data []byte) error {
	return json.Unmarshal(data, jp.data)
}
