package mapper

import (
	`fmt`
	`github.com/exklamationmark/glog`
	. `gopkg.in/check.v1`
	`time`
)

func getTime(timeStr string) time.Time {
	loc, _ := time.LoadLocation(`Singapore`)
	parsedTime, err := time.ParseInLocation(`2006-01-02 15:04:05-07`, timeStr, loc)
	if err != nil {
		glog.Fatal(fmt.Sprintf(`cannot parse time, time="%s", err=%v`, timeStr, err))
	}
	return parsedTime
}

type testEntry struct {
	target     interface{}
	comparator Checker
	value      interface{}
}

// perform tests based on table, to dry up similar test cases
// 2 transformers can be given, first one for target, 2nd for value
func tableCheck(c *C, tests []testEntry, transformers ...func(interface{}) interface{}) {
	var targetT, valueT func(interface{}) interface{}
	//var targetT, valueT testTransformer
	var target, value interface{}
	if len(transformers) > 0 {
		targetT = transformers[0] //.(testTransformer)
	}
	if len(transformers) > 1 {
		valueT = transformers[1] //.(testTransformer)
	}

	for _, test := range tests {
		target = test.target
		if targetT != nil {
			target = targetT(target)
		}
		value = test.value
		if valueT != nil {
			value = valueT(value)
		}

		c.Assert(target, test.comparator, value)
	}
}

// transform the target into actual value from a Record map
func recordCheck(rec Record, tests []testEntry, c *C) {
	for _, test := range tests {
		c.Assert(rec[test.target.(string)], test.comparator, test.value)
	}
}

// a wrapper, so don't have to do error check
func exec(query string, args ...interface{}) {
	_, err := dbconnection.Exec(query, args...)
	if err != nil {
		glog.Fatal(fmt.Sprintf(`error running query, query= %v; args= %v; err= %v`, query, args, err))
	}
}

// fixture: create a fake table to test, something similar to a user record
func createTestTables() {
	exec(`DROP TABLE IF EXISTS t_users, t_roles, t_user_roles`)
	exec(`CREATE TABLE t_users (
		id character varying(15) NOT NULL PRIMARY KEY,
		email character varying(255),
		age integer NOT NULL,
		active boolean NOT NULL,
		email_verified boolean,
		no_of_licenses integer,
		last_payment_at timestamp with time zone,
		created_at timestamp with time zone NOT NULL
	)`)
	exec(`CREATE TABLE t_roles (
		id character varying(15) NOT NULL PRIMARY KEY,
		name character varying(255) NOT NULL,
		required_karma integer NOT NULL
	)`)
	exec(`CREATE TABLE t_user_roles (
		id character varying(15) NOT NULL PRIMARY KEY,
		user_id character varying(15) NOT NULL,
		role_id character varying(15) NOT NULL
	)`)
	exec(`CREATE UNIQUE INDEX ON t_user_roles (user_id, role_id)`)
	exec(`TRUNCATE TABLE t_users, t_roles, t_user_roles`)
}

func testQuery(c *C, q *Query, queryType int, query string, selectFields []string, args []interface{}) {
	var tests = []testEntry{
		{q.queryType, Equals, queryType},
		{q.query, Equals, query},
		{q.selectFields, DeepEquals, selectFields},
		{q.args, DeepEquals, args},
	}
	tableCheck(c, tests)
}
