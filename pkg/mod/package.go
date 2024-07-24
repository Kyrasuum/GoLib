package mod

import (
	"gorm.io/gorm"

	"library/pkg/io"
)

// for reference on available hooks for gorm
type model interface {
	AfterFind(*gorm.DB) error
	AfterCreate(*gorm.DB) error
	BeforeCreate(*gorm.DB) error
	AfterSave(*gorm.DB) error
	BeforeSave(*gorm.DB) error
	AfterUpdate(*gorm.DB) error
	BeforeUpdate(*gorm.DB) error
	BeforeDelete(*gorm.DB) error
	AfterDelete(*gorm.DB) error
}

type recursFunction func(string, ...interface{}) (interface{}, error)
type moduleFunction func(recursFunction, ...interface{}) (interface{}, error)

type module map[string]interface{}

var (
	modules []module
)

func LoadModule(mod interface{}) {
	switch modtyped := mod.(type) {
	case module:
		modules = append(modules, modtyped)
	default:
		modmap, err := io.StructToMap(modtyped)
		if err != nil {
			return
		}
		modules = append(modules, modmap)
	}
}

func CallModuleFunction(function string, params ...interface{}) (interface{}, error) {
	for _, mod := range modules {
		if modfunc, ok := mod[function]; ok {
			ret, err := modfunc.(moduleFunction)(CallModuleFunction, params...)
			if err != nil || ret != nil {
				return ret, err
			}
		}
	}
	return nil, nil
}
