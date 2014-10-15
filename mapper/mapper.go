// Pacakge mapper provide a cleaner query interface for postgres
// We start by connecting mapper to an existing db connection, then register tables in there
// After that query can be constructed and the result will be put into a map
package mapper

import (
	`fmt`
	`github.com/exklamationmark/glog`
)

const (
	initColCount      = 15
	initTableCount    = 30
	initTotalColCount = initColCount * initTableCount
)

const (
	stringType = iota
	nullStringType
	int64Type
	nullInt64Type
	boolType
	nullBoolType
	timeType
	nullTimeType
	invalidType
)

var (
	schemaQuery = `SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '%s'`

	invalidTypeErr      = `invalid sql data type, got "%v", expected one of ("character varying", "text", "integer", "boolean", "timestamp with time zone", "timestamp without time zone")`
	invalidNullableErr  = `invalid value for nullable, got "%v", expected one of ("YES", "NO")`
	cannotLoadSchemaErr = `cannot load schema, error= %v`

	//tables = make(map[string]*Table, initTableCount)
	columns = make(map[string]int, initTotalColCount)
)

// Register query the db for a table's schema and store them for later use
func Register(tbName string) {
	rows, err := dbconnection.Query(fmt.Sprintf(schemaQuery, tbName))
	defer rows.Close()
	if err != nil {
		glog.Fatal(fmt.Errorf(cannotLoadSchemaErr, err))
	}

	defer rows.Close()
	var colName, dataType, nullable string
	for rows.Next() {
		if err := rows.Scan(&colName, &dataType, &nullable); err != nil {
			glog.Fatal(fmt.Errorf(cannotLoadSchemaErr, err))
			continue
		}

		colType, err := toType(dataType, nullable)
		if err != nil {
			glog.Fatal(fmt.Errorf(cannotLoadSchemaErr, err))
		}

		columns[tbName+`.`+colName] = colType
	}
}

// toType returns the corresponding Golang type for a sql data_type
func toType(dataType, nullable string) (int, error) {
	if nullable != `NO` && nullable != `YES` {
		return invalidType, fmt.Errorf(invalidNullableErr, nullable)
	}

	switch dataType {
	case `character varying`, `text`, `inet`:
		if nullable == `NO` {
			return stringType, nil
		}
		return nullStringType, nil
	case `integer`:
		if nullable == `NO` {
			return int64Type, nil
		}
		return nullInt64Type, nil
	case `boolean`:
		if nullable == `NO` {
			return boolType, nil
		}
		return nullBoolType, nil
	case `timestamp with time zone`, `timestamp without time zone`:
		if nullable == `NO` {
			return timeType, nil
		}
		return nullTimeType, nil
	}

	return invalidType, fmt.Errorf(invalidTypeErr, dataType)
}
