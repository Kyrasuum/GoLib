package db

import (
	"gorm.io/gorm"
)

// generic pager
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 0 {
			page = 0
		}

		if pageSize <= 0 {
			pageSize = 1
		}

		offset := page * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
