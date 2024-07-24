package io

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

type YamlParser parser

func (yp YamlParser) Save(file string) (err error) {
	if len(file) < 1 {
		return fmt.Errorf("invalid json file\n%s", err)
	}

	// convert object into a json bytelike object
	yamldata, err := yp.ToYaml()
	if err != nil {
		return fmt.Errorf("failed to marshal object\n%s", err)
	}

	// write json bytelike object to file
	err = ioutil.WriteFile(file, yamldata, 0664)
	if err != nil {
		return fmt.Errorf("failed to write to json file\n%s", err)
	}

	return err
}

func (yp YamlParser) Load(file string, structobj interface{}) (err error) {
	data, err := ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error opening yaml/yml file\n%s", err)
	}

	// convert yml into a map[string]interface{}
	err = yp.FromYaml(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml/yml file\n%s", err)
	}

	// convert map into a json object
	// this allows shared field processing with json load
	jsonObj, err := json.Marshal(yp.data)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml/yml file\n%s", err)
	}

	// converts json object into a struct
	err = json.Unmarshal(jsonObj, structobj)
	if err != nil {
		return fmt.Errorf("failed to parse yaml/yml file\n%s", err)
	}
	yp.data = structobj

	return err
}

func (yp YamlParser) ToMap() (map[string]interface{}, error) {
	temp, err := yp.ToYaml()
	if err != nil {
		return nil, err
	}
	return yp.data.(map[string]interface{}), yp.FromYaml(temp)
}

func (yp YamlParser) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, yp.data)
}

func (yp YamlParser) ToYaml() ([]byte, error) {
	return yaml.Marshal(yp.data)
}
