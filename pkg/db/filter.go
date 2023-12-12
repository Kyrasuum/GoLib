package db

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// perform string search in database
func StringSearchDB(table string, col string, fields []string, req *gorm.DB) *gorm.DB {
	if table != "" {
		if table[0] == '!' {
			return req.Where(fmt.Sprintf("NOT CAST (%s.%s as text) ILIKE ANY(ARRAY['%s'])", table[1:], col, strings.ReplaceAll(strings.Join(fields, "','"), "*", "%")))
		} else {
			return req.Where(fmt.Sprintf("CAST (%s.%s as text) ILIKE ANY(ARRAY['%s'])", table, col, strings.ReplaceAll(strings.Join(fields, "','"), "*", "%")))
		}
	} else {
		if col[0] == '!' {
			return req.Where(fmt.Sprintf("NOT CAST (%s as text) ILIKE ANY(ARRAY['%s'])", col[1:], strings.ReplaceAll(strings.Join(fields, "','"), "*", "%")))
		} else {
			return req.Where(fmt.Sprintf("CAST (%s as text) ILIKE ANY(ARRAY['%s'])", col, strings.ReplaceAll(strings.Join(fields, "','"), "*", "%")))
		}
	}
}

// perform date search in database
func DateSearchDB(table string, col string, fields []string, req *gorm.DB) *gorm.DB {
	search := "0 = 1"
	for _, field := range fields {
		dates := strings.Split(strings.ReplaceAll(field, "*", "%"), "->")
		if len(dates) > 1 {
			if table != "" {
				if table[0] == '!' {
					search = fmt.Sprintf("%s OR NOT %s.%s BETWEEN '%s' AND '%s'", search, table[1:], col, dates[0], dates[1])
				} else {
					search = fmt.Sprintf("%s OR %s.%s BETWEEN '%s' AND '%s'", search, table, col, dates[0], dates[1])
				}
			} else {
				if col[0] == '!' {
					search = fmt.Sprintf("%s OR NOT %s BETWEEN '%s' AND '%s'", search, col[1:], dates[0], dates[1])
				} else {
					search = fmt.Sprintf("%s OR %s BETWEEN '%s' AND '%s'", search, col, dates[0], dates[1])
				}
			}
		}
	}
	return req.Where(search)
}

// apply generic filtering in database
func ApplyFiltersDB(table string, submap map[string][]string, req *gorm.DB) *gorm.DB {
	for col, fields := range submap {
		table := table
		if table == "" {
			split := strings.Split(col, ".")
			table = split[0]
			col = split[1]
		}
		switch col {
		case "created_at", "created", "updated_at", "date":
			req = DateSearchDB(table, col, fields, req)
		default:
			req = StringSearchDB(table, col, fields, req)
		}
	}
	return req
}
