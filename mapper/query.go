package mapper

// Query corresponds to an actual query to be made
type Query struct {
	query        string
	args         []interface{}
	selectFields []string
	queryType    int
}

const (
	SelectQuery = iota
	InsertQuery
	BulkInsertQuery
	UpdateQuery
	DeleteQuery
	TruncateQuery
)
