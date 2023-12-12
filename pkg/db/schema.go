package pkg

import (
	"fmt"
	"os/exec"
	"strings"

	"library/pkg/io"
)

func Schema(DataPath string, DataSet string, DebugSQL bool) (err error) {
	out, err := exec.Command("ls", DataPath).Output()
	if err != nil {
		return fmt.Errorf("could not open data directory\n%s", err)
	}
	dirs := strings.Split(fmt.Sprint(string(out)), "\n")
	obj := map[string]interface{}{}

	for _, dir := range dirs {
		if DataSet != "all" && string(dir) != DataSet {
			continue
		}

		err = LoadSchemaYml(DataPath, fmt.Sprint(string(dir)), "", &obj, DebugSQL)
		if err != nil {
			return err
		}
		err = LoadSchemaJson(DataPath, fmt.Sprint(string(dir)), "", &obj, DebugSQL)
		if err != nil {
			return err
		}
		err = LoadSchemaYaml(DataPath, fmt.Sprint(string(dir)), "", &obj, DebugSQL)
		if err != nil {
			return err
		}
	}

	for key, value := range obj {
		fmt.Printf("\n%s:\n", key)
		subbytes, err := io.UnescapeJson(value)
		if err == nil {
			fmt.Printf("%s\n", subbytes)
		} else {
			fmt.Printf("%+v\n", err)
		}
		fmt.Printf("\n")
	}

	return err
}

func LoadSchemaJson(DataPath string, Repo string, Prefix string, obj *map[string]interface{}, DebugSQL bool) (err error) {
	out, _ := exec.Command("/usr/bin/find", DataPath+Repo, "-wholename", Prefix+"**/*.json").CombinedOutput()
	files := strings.Split(fmt.Sprint(string(out)), "\n")
	for _, file := range files {
		if len(file) < 1 {
			continue
		}

		// load the file into structs
		var subobj interface{}
		err := io.LoadJsonFile(file, &subobj)
		if err != nil {
			fmt.Printf("Failed to load %s\n", file)
			return err
		}

		//remove contents
		(*obj)[Repo+"/"+Prefix] = SchemaProc(file, subobj, (*obj)[Repo+"/"+Prefix])
	}
	return err
}

func LoadSchemaYml(DataPath string, Repo string, Prefix string, obj *map[string]interface{}, DebugSQL bool) (err error) {
	out, _ := exec.Command("/usr/bin/find", DataPath+Repo, "-wholename", Prefix+"**/*.yml").CombinedOutput()
	files := strings.Split(fmt.Sprint(string(out)), "\n")
	for _, file := range files {
		if len(file) < 1 {
			continue
		}

		// load file into a struct
		var subobj interface{}
		err := io.LoadYmlFile(file, &subobj)
		if err != nil {
			fmt.Printf("Failed to load %s\n", file)
			return err
		}

		//remove contents
		(*obj)[Repo+"/"+Prefix] = SchemaProc(file, subobj, (*obj)[Repo+"/"+Prefix])
	}
	return err
}

func LoadSchemaYaml(DataPath string, Repo string, Prefix string, obj *map[string]interface{}, DebugSQL bool) (err error) {
	out, _ := exec.Command("/usr/bin/find", DataPath+Repo, "-wholename", Prefix+"**/*.yaml").CombinedOutput()
	files := strings.Split(fmt.Sprint(string(out)), "\n")
	for _, file := range files {
		if len(file) < 1 {
			continue
		}
		// load file into a struct
		var subobj interface{}
		err := io.LoadYamlFile(file, &subobj)
		if err != nil {
			fmt.Printf("Failed to load %s\n", file)
			return err
		}

		//remove contents
		(*obj)[Repo+"/"+Prefix] = SchemaProc(file, subobj, (*obj)[Repo+"/"+Prefix])
	}
	return err
}

func SchemaProc(file string, obj interface{}, schema interface{}) interface{} {
	var subschema interface{}

	defer func() {
		if recover() != nil {
			fmt.Printf("\nFailed to load: %s\n", file)
			fmt.Printf("\tobj: %+v\t\t%T\n", obj, obj)
			fmt.Printf("\tschema: <%+v>\t\t%T\n\n", schema, schema)
		}
	}()

	switch obj.(type) {
	case string:
		subschema = ""
	case []interface{}:
		if schema != nil && schema != "" && len(schema.([]interface{})) > 0 {
			subschema = schema
		} else {
			subschema = []interface{}{}
			schema = []interface{}{}
		}
		for key := range obj.([]interface{}) {
			if len(schema.([]interface{})) > 0 {
				subschema.([]interface{})[0] = SchemaProc(file, obj.([]interface{})[key], schema.([]interface{})[0])
			} else {
				var temp interface{}
				subschema = append(schema.([]interface{}), SchemaProc(file, obj.([]interface{})[key], temp))
			}
			schema = subschema
		}
	case map[string]interface{}:
		if schema != nil && schema != "" {
			subschema = schema
		} else {
			subschema = map[string]interface{}{}
			schema = map[string]interface{}{}
		}
		for key := range obj.(map[string]interface{}) {
			if _, ok := schema.(map[string]interface{})[key]; ok {
				subschema.(map[string]interface{})[key] = SchemaProc(file, (obj.(map[string]interface{})[key]), schema.(map[string]interface{})[key])
			} else {
				var temp interface{}
				subschema.(map[string]interface{})[key] = SchemaProc(file, (obj.(map[string]interface{})[key]), temp)
			}
			schema.(map[string]interface{})[key] = subschema.(map[string]interface{})[key]
		}
	case float64:
		subschema = 0.0
	case bool:
		subschema = true
	case nil:
		subschema = schema
	default:
		fmt.Printf("Unknown subobj type: %T\n", obj)
	}
	return subschema
}
