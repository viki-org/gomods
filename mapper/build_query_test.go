package mapper

import (
	. `gopkg.in/check.v1`
)

type BuilderTS struct {
	query Query
}

func init() {
	Suite(&BuilderTS{})
}

// use a table here to minize copying code
var scenarios = []struct {
	q            *Query
	queryType    int
	query        string
	selectFields []string
	args         []interface{}
}{
	{
		Select(`t_users.id`, `t_users.email`, `t_users.email_verified`).From(`t_users`),
		SelectQuery,
		`SELECT t_users.id, t_users.email, t_users.email_verified FROM t_users`,
		[]string{`t_users.id`, `t_users.email`, `t_users.email_verified`},
		nil,
	},
	{
		Select(`t_roles.name`, `t_roles.required_karma`).FromJoin(InnerJoin, `t_user_roles`, `t_roles`, `t_user_roles.role_id = t_roles.id`),
		SelectQuery,
		`SELECT t_roles.name, t_roles.required_karma FROM t_user_roles INNER JOIN t_roles ON t_user_roles.role_id = t_roles.id`,
		[]string{`t_roles.name`, `t_roles.required_karma`},
		nil,
	},
	{
		Select(`t_users.id`, `t_users.email`, `t_users.email_verified`).From(`t_users`).Where(`t_users.id = ? AND t_users.email_verified = ?`, `10u`, false),
		SelectQuery,
		`SELECT t_users.id, t_users.email, t_users.email_verified FROM t_users WHERE t_users.id = $1 AND t_users.email_verified = $2`,
		[]string{`t_users.id`, `t_users.email`, `t_users.email_verified`},
		[]interface{}{`10u`, false},
	},
	{
		Select(`t_users.id`, `t_users.email`, `t_users.email_verified`).From(`t_users`).Order(`t_users.id`, Asc).Limit(20),
		SelectQuery,
		`SELECT t_users.id, t_users.email, t_users.email_verified FROM t_users ORDER BY t_users.id ASC LIMIT 20`,
		[]string{`t_users.id`, `t_users.email`, `t_users.email_verified`},
		nil,
	},
	{
		Insert(`t_roles`, `id, name, required_karma`, `1r`, `Code monkey`, 100),
		InsertQuery,
		`INSERT INTO t_roles (id, name, required_karma) VALUES ($1, $2, $3)`,
		nil,
		[]interface{}{`1r`, `Code monkey`, 100},
	},
	{
		Update(`t_roles`, `name = ?, required_karma = ?`, `Code kingkong`, 500),
		UpdateQuery,
		`UPDATE t_roles SET name = $1, required_karma = $2`,
		nil,
		[]interface{}{`Code kingkong`, 500},
	},
	{
		Update(`t_roles`, `name = ?, required_karma = ?`, `Bug eagle`, 1000).Where(`id = ?`, `1r`),
		UpdateQuery,
		`UPDATE t_roles SET name = $1, required_karma = $2 WHERE id = $3`,
		nil,
		[]interface{}{`Bug eagle`, 1000, `1r`},
	},
	{
		Delete(`t_roles`),
		DeleteQuery,
		`DELETE FROM t_roles`,
		nil,
		nil,
	},
	{
		Delete(`t_roles`).Where(`id = ?`, `1r`),
		DeleteQuery,
		`DELETE FROM t_roles WHERE id = $1`,
		nil,
		[]interface{}{`1r`},
	},
	{
		Truncate(`t_roles`, `t_users`),
		TruncateQuery,
		`TRUNCATE t_roles, t_users`,
		nil,
		nil,
	},
}

func (s *BuilderTS) TestQueryBuilder(c *C) {
	for _, test := range scenarios {
		testQuery(c, test.q, test.queryType, test.query, test.selectFields, test.args)
	}
}
