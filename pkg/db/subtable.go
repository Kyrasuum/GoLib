package db

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"library/pkg"
)

// applies a function to each nested table in passed object
type applyFunc func(req *gorm.DB, table string, field string, obj interface{}, params ...interface{}) (*gorm.DB, error)

func ApplyFuncSubtables(req *gorm.DB, record interface{}, apply applyFunc, params ...interface{}) (subreq *gorm.DB, err error) {
	t := reflect.TypeOf(record)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		//detect if field is a nested table
		val := reflect.ValueOf(record).Field(i).Interface()
		vt := fmt.Sprintf("*%T", val)
		tab, err := pkg.CallModuleFunction("GetTable", vt)
		if err == nil {
			subreq, err = apply(req, tab.(string), field.Name, val, params...)
			if err == nil {
				return subreq, err
			}
		}
	}
	return subreq, nil
}
