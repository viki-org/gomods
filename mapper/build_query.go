package mapper

import (
	`database/sql`
	`fmt`
	`strings`
)

var (
	selectTemplate   = `SELECT %s`
	fromTemplate     = `%s FROM %s`
	fromJoinTemplate = `%s FROM %s %s %s ON %s`
	whereTemplate    = `%s WHERE %s`
	argTemplate      = `$%v`
	limitTemplate    = `%s LIMIT %d`
	orderTemplate    = `%s ORDER BY %s %s`
	insertTemplate   = `INSERT INTO %s (%s) VALUES (%s)`
	updateTemplate   = `UPDATE %s SET %s`
	deleteTemplate   = `DELETE FROM %s`
	truncateTemplate = `TRUNCATE %s`
)

var dbconnection *sql.DB

// Connect register a db connection to the module. Must be called to init mapper
func Connect(conn *sql.DB) {
	dbconnection = conn
}

// Select starts the creation of a select query
func Select(fields ...string) *Query {
	return &Query{
		queryType:    SelectQuery,
		query:        fmt.Sprintf(selectTemplate, strings.Join(fields, `, `)),
		selectFields: fields,
	}
}

// From indicates a table for the query
func (q *Query) From(table string) *Query {
	q.query = fmt.Sprintf(fromTemplate, q.query, table)
	return q
}

const (
	InnerJoin = iota
)

var (
	joinWords = map[int]string{
		InnerJoin: `INNER JOIN`,
	}
)

// FromJoin indicates a joint of tables as the source, for now only take cares of 2 table join
func (q *Query) FromJoin(joinType int, first, second, conditions string) *Query {
	q.query = fmt.Sprintf(fromJoinTemplate, q.query, first, joinWords[joinType], second, conditions)
	return q
}

const (
	placeHolder = `?`
)

// Where construct the where clause of the query and add arugments for it
// use ? for place holders
// assume no of `?` in conditions & no of args is the same
func (q *Query) Where(conditions string, args ...interface{}) *Query {
	parts := strings.Split(conditions, placeHolder)
	final := make([]string, 0, len(parts)*2)
	start := len(q.args) + 1
	for offset := range args {
		final = append(final, parts[offset], fmt.Sprintf(argTemplate, start+offset))
	}
	q.query = fmt.Sprintf(whereTemplate, q.query, strings.Join(final, ``))
	q.args = append(q.args, args...)
	return q
}

const (
	Asc = iota
	Desc
)

var (
	orderWords = map[int]string{
		Asc:  `ASC`,
		Desc: `DESC`,
	}
)

// Order sorts the returning rows
func (q *Query) Order(field string, orderType int) *Query {
	q.query = fmt.Sprintf(orderTemplate, q.query, field, orderWords[orderType])
	return q
}

// Limit constrain the no of rows to return, and hence no of lookup
func (q *Query) Limit(limit int) *Query {
	q.query = fmt.Sprintf(limitTemplate, q.query, limit)
	return q
}

// Insert starts an insert query
// assume no of fields == no of args
func Insert(table, fields string, args ...interface{}) *Query {
	argsStr := make([]string, 0, len(args))
	for index := range args {
		argsStr = append(argsStr, fmt.Sprintf(argTemplate, index+1))
	}
	return &Query{
		queryType: InsertQuery,
		query:     fmt.Sprintf(insertTemplate, table, fields, strings.Join(argsStr, `, `)),
		args:      args,
	}
}

// Update starts an update query
// assume no of fields == no of args
func Update(table, fields string, args ...interface{}) *Query {
	parts := strings.Split(fields, placeHolder)
	final := make([]string, 0, len(parts)*2)
	for index := range args {
		final = append(final, parts[index], fmt.Sprintf(argTemplate, index+1))
	}
	return &Query{
		queryType: UpdateQuery,
		query:     fmt.Sprintf(updateTemplate, table, strings.Join(final, ``)),
		args:      args,
	}
}

// Delete starts a delete query
// be careful and add a where clause, or you will truncate the whole table
func Delete(table string) *Query {
	return &Query{
		queryType: DeleteQuery,
		query:     fmt.Sprintf(deleteTemplate, table),
	}
}

// Truncate starts a truncate query
func Truncate(tables ...string) *Query {
	return &Query{
		queryType: TruncateQuery,
		query:     fmt.Sprintf(truncateTemplate, strings.Join(tables, `, `)),
	}
}

// Run executes a query
func (q *Query) Run() ([]Record, error) {
	return Exec(q)
}
