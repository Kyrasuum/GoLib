package db

import (
	"gorm.io/gorm"
)

// applies sort parameters to a gorm query
func SortResults(sorts []interface{}, req *gorm.DB) *gorm.DB {
	for _, sort := range sorts {
		switch sort.(string)[0:1] {
		case "+":
			req = req.Order(sort.(string)[1:] + " asc")
		case "-":
			req = req.Order(sort.(string)[1:] + " desc")
		default:
			req = req.Order(sort.(string)[0:])
		}
	}
	return req
}
