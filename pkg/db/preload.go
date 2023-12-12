package db

import (
	"gorm.io/gorm"
)

func PreloadStruct(req *gorm.DB, table string, field string, obj interface{}, params ...interface{}) (*gorm.DB, error) {
	prefix := params[0].(string)
	depth := params[1].(int)
	limit := params[2].(int)
	if depth < limit || limit < 0 {
		req = req.Preload(prefix + field)
		return ApplyFuncSubtables(req, obj, PreloadStruct, "", 0, limit)
	}
	return req, nil
}
func GreedyPreload(req *gorm.DB, record interface{}, limit int) (*gorm.DB, error) {
	return ApplyFuncSubtables(req, record, PreloadStruct, "", 0, limit)
}
