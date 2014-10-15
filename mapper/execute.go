package mapper

import (
	`database/sql`
	`fmt`
	`github.com/exklamationmark/glog`
	`github.com/lib/pq`
	`time`
)

const (
	cannotRunQueryErr = `query "%s" failed to run, err=%v`
	rowScanErr        = `scanning row failed, rows=%v, err=%v`

	initResultsCount = 10
)

// Record is a map of columns / values returned by a query (SelectQuery)
type Record map[string]interface{}

// Exec run a query and extract results as a map
func Exec(query *Query) ([]Record, error) {
	rows, err := dbconnection.Query(query.query, query.args...)
	if err != nil {
		glog.Error(fmt.Printf(cannotRunQueryErr, query.query, err))
		return nil, err
	}
	defer rows.Close()

	// assuming other queries doesn't return values, only Select does
	if query.queryType != SelectQuery {
		return nil, nil
	}

	results := make([]Record, 0, initResultsCount)
	placeholders := createPlaceholders(query.selectFields)

	for rows.Next() {

		if err := rows.Scan(placeholders...); err != nil {
			glog.Error(fmt.Sprintf(rowScanErr, rows, err))
			return nil, err
		}

		// placeholders should contain data in order of fields in selectFields
		record := make(Record, len(query.selectFields))
		for i := 0; i < len(query.selectFields); i++ {
			colType := columns[query.selectFields[i]]
			switch colType {
			case stringType:
				record[query.selectFields[i]] = *(placeholders[i].(*string))
			case nullStringType:
				record[query.selectFields[i]] = *(placeholders[i].(*sql.NullString))
			case int64Type:
				record[query.selectFields[i]] = *(placeholders[i].(*int64))
			case boolType:
				record[query.selectFields[i]] = *(placeholders[i].(*bool))
			case nullBoolType:
				record[query.selectFields[i]] = *(placeholders[i].(*sql.NullBool))
			case nullInt64Type:
				record[query.selectFields[i]] = *(placeholders[i].(*sql.NullInt64))
			case nullTimeType:
				record[query.selectFields[i]] = *(placeholders[i].(*pq.NullTime))
			case timeType:
				record[query.selectFields[i]] = *(placeholders[i].(*time.Time))
			default:
				return nil, fmt.Errorf(`unknown column type`)
			}
		}
		results = append(results, record)
	}

	return results, nil
}

// createPlaceholder generate a slice of pointers to hold data in select query
func createPlaceholders(fields []string) []interface{} {
	placeholders := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		fieldType := columns[field]
		switch fieldType {
		case stringType:
			placeholders[i] = new(string)
		case nullStringType:
			placeholders[i] = new(sql.NullString)
		case int64Type:
			placeholders[i] = new(int64)
		case boolType:
			placeholders[i] = new(bool)
		case nullBoolType:
			placeholders[i] = new(sql.NullBool)
		case nullInt64Type:
			placeholders[i] = new(sql.NullInt64)
		case nullTimeType:
			placeholders[i] = new(pq.NullTime)
		case timeType:
			placeholders[i] = new(time.Time)
		}
	}
	return placeholders
}
