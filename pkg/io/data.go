package io

import (
	"bytes"
	"encoding/json"

	"sigs.k8s.io/yaml"
)

func UnescapeJson(message interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(message)

	return buffer.Bytes(), err
}

// convert yml into a interface{}
func YamlToMap(data []byte) (ret interface{}, err error) {
	err = yaml.Unmarshal(data, &ret)
	return ret, err
}

// convert a struct to a map
func StructToMap(data interface{}) (map[string]interface{}, error) {
	temp, err := PackJson(data)
	if err != nil {
		return nil, err
	}
	ret, err := UnpackJson(temp, map[string]interface{}{})
	return ret.(map[string]interface{}), err
}

// convert interface{} into a json string object
func PackJson(data interface{}) (ret []byte, err error) {
	ret, err = json.Marshal(data)
	return ret, err
}

// converts json object into a struct
func UnpackJson(data []byte, dest interface{}) (interface{}, error) {
	err := json.Unmarshal(data, dest)
	return dest, err
}
