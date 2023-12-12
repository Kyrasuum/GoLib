package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// retrieve gorm primary keys of a struct
func GetPrimaryKeys(record interface{}) (ret interface{}, err error) {
	ret = reflect.New(reflect.TypeOf(record)).Interface()
	t := reflect.TypeOf(record)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("gorm")
		if strings.Contains(tag, "primaryKey") {
			indirect := reflect.Indirect(reflect.ValueOf(record))
			val := indirect.FieldByName(field.Name)
			reflect.ValueOf(ret).Elem().FieldByName(field.Name).Set(val)
		}
	}

	return ret, nil
}

// retrieve table columns
func GetTableColumns(db *gorm.DB, table interface{}, DebugSQL bool) (cols []string, err error) {
	if db == nil {
		return nil, errors.New("database not connected")
	}

	req := ConditionalDebug(db.Session(&gorm.Session{PrepareStmt: true}), DebugSQL)
	if err = req.Raw(fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s'", "public", table)).Find(&cols).Error; err != nil {
		return nil, err
	}

	return cols, nil
}

// gets all table names contained in database
func GetTableNames(db *gorm.DB, DebugSQL bool) (tables []string, err error) {
	if db == nil {
		return nil, errors.New("database not connected")
	}

	req := ConditionalDebug(db.Session(&gorm.Session{PrepareStmt: true}).Table("information_schema.tables"), DebugSQL)
	if err := req.Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
}

// empty database
func FlushDatabase(db *gorm.DB) error {
	fmt.Printf("Flushing Database\n")
	db.Exec("DO $$\n" +
		"DECLARE r RECORD;\n" +
		"BEGIN FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema())\n" +
		"LOOP\n" +
		"EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';\n" +
		"END LOOP;\n" +
		"END $$;")
	fmt.Printf("Flushing Complete\n")

	return nil
}
