package db

import (
	"gorm.io/gorm"
)

// generic rollback function
func DBrollback(tx *gorm.DB, err *error) {
	if r := recover(); r != nil {
		tx.Rollback()
		*err = r.(error)
	}
}

// conditionally print sql debugging
func ConditionalDebug(req *gorm.DB, DebugSQL bool) *gorm.DB {
	if DebugSQL {
		return req.Debug()
	} else {
		return req
	}
}
