package io

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

func SaveJsonFile(file string, subobj interface{}) (err error) {
	if len(file) < 1 {
		return fmt.Errorf("invalid json file\n%s", err)
	}

	// convert object into a json bytelike object
	jsonObj, err := UnescapeJson(subobj)
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

func LoadJsonFile(file string, subobj interface{}) (err error) {
	if len(file) < 1 {
		return fmt.Errorf("invalid json file\n%s", err)
	}

	// Read in the files
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to open json file\n%s", err)
	}

	// Parse the json into structs
	err = json.Unmarshal(data, subobj)
	if err != nil {
		return fmt.Errorf("failed to parse json file\n%s", err)
	}

	return err
}

func LoadYmlFile(file string, subobj interface{}) (err error) {
	if len(file) < 1 {
		return fmt.Errorf("invalid yml file\n%s", err)
	}

	// Read in the files
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to open yml file\n%s", err)
	}

	// convert yml into a map[string]interface{}
	var yamlObj interface{}
	err = yaml.Unmarshal(data, &yamlObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yml file\n%s", err)
	}

	// convert map into a json object
	jsonObj, err := json.Marshal(yamlObj)
	if err != nil {
		return fmt.Errorf("failed to marshal yml file\n%s", err)
	}

	// converts json object into a struct
	err = json.Unmarshal(jsonObj, subobj)
	if err != nil {
		return fmt.Errorf("failed to parse yml file\n%s", err)
	}

	return err
}

func LoadYamlFile(file string, subobj interface{}) (err error) {
	if len(file) < 1 {
		return fmt.Errorf("invalid yaml file\n%s", err)
	}

	// Read in the files
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to open yaml file\n%s", err)
	}

	// convert yaml into a map[string]interface{}
	var yamlObj interface{}
	err = yaml.Unmarshal(data, &yamlObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml file\n%s", err)
	}

	// convert map into a json object
	jsonObj, err := json.Marshal(yamlObj)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml file\n%s", err)
	}

	// converts json object into a struct
	err = json.Unmarshal(jsonObj, subobj)
	if err != nil {
		return fmt.Errorf("failed to parse yaml file\n%s", err)
	}

	return err
}
