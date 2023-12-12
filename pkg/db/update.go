package db

import (
	"gorm.io/gorm"
)

// apply update on sub table
func UpdateSubTable(req *gorm.DB, table string, field string, obj interface{}, params ...interface{}) (*gorm.DB, error) {
	//attempt to update nested tables
	req = req.Preload(field)
	asc := req.Association(field)
	asc.Replace(obj)
	if err := asc.Error; err != nil {
		return req, err
	}
	return req, nil
}

// cascade update sub tables
func UpdatedSubTables(req *gorm.DB, record interface{}) (*gorm.DB, error) {
	return ApplyFuncSubtables(req, record, UpdateSubTable)
}
