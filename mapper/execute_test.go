package mapper

import (
	`database/sql`
	`github.com/exklamationmark/glog`
	`github.com/lib/pq`
	. `github.com/viki-org/gomods/sqlcheckers`
	. `gopkg.in/check.v1`
)

type SelectExecTS struct {
	query *Query
}

type InsertExecTS struct {
	query *Query
}

type BulkInsertExecTS struct {
	query *Query
}

type UpdateExecTS struct {
	query *Query
}

type DeleteExecTS struct {
	query *Query
}

func init() {
	conn, err := sql.Open(`postgres`, `host=localhost port=5432 sslmode=disable dbname=users_test user=postgres password=password`)
	if err != nil {
		glog.Fatal(`cannot connect to postgres, err=`, err)
	}
	Connect(conn)

	Suite(&SelectExecTS{})
	Suite(&InsertExecTS{})
	Suite(&BulkInsertExecTS{})
	Suite(&UpdateExecTS{})
	Suite(&DeleteExecTS{})
}

var (
	testTime     = getTime(`2014-06-18 09:00:00+00`)
	sampleSelect = &Query{
		queryType:    SelectQuery,
		query:        `SELECT t_users.id, t_users.email, t_users.age, t_users.active, t_users.email_verified, t_users.no_of_licenses, t_users.last_payment_at, t_users.created_at FROM t_users`,
		selectFields: []string{`t_users.id`, `t_users.email`, `t_users.age`, `t_users.active`, `t_users.email_verified`, `t_users.no_of_licenses`, `t_users.last_payment_at`, `t_users.created_at`},
	}
	sampleInsert = &Query{
		queryType: InsertQuery,
		query:     `INSERT INTO t_users (id, email, age, active, email_verified, no_of_licenses, last_payment_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		args:      []interface{}{`2u`, `darth@vader.com`, 40, true, true, 0, nil, testTime},
	}
)

func (s *SelectExecTS) SetUpTest(c *C) {
	createTestTables()
	Register(`t_users`)
	exec(`INSERT INTO t_users (id, email, age, active, email_verified, no_of_licenses, last_payment_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, `1u`, `user@test.com`, 20, false, false, 0, nil, testTime)

	s.query = sampleSelect
}

func (s *SelectExecTS) TestSelectExec(c *C) {
	data, err := Exec(s.query)

	c.Assert(err, IsNil)
	c.Assert(len(data), Equals, 1)

	var tests = []testEntry{
		{`t_users.id`, Equals, `1u`},
		{`t_users.email`, SQLEquals, sql.NullString{Valid: true, String: `user@test.com`}},
		{`t_users.age`, Equals, int64(20)},
		{`t_users.active`, Equals, false},
		{`t_users.email_verified`, SQLEquals, sql.NullBool{Valid: true, Bool: false}},
		{`t_users.no_of_licenses`, SQLEquals, sql.NullInt64{Valid: true, Int64: int64(0)}},
		{`t_users.created_at`, SQLEquals, testTime},
		{`t_users.last_payment_at`, SQLEquals, pq.NullTime{Valid: false, Time: testTime}},
	}
	recordCheck(data[0], tests, c)
}

func (s *InsertExecTS) SetUpTest(c *C) {
	createTestTables()
	Register(`t_users`)

	s.query = sampleInsert
}

func (s *InsertExecTS) TestInsertExec(c *C) {
	data, err := Exec(s.query)

	c.Assert(err, IsNil)
	c.Assert(data, IsNil)

	data, err = Exec(sampleSelect)
	var tests = []testEntry{
		{`t_users.id`, Equals, `2u`},
		{`t_users.email`, SQLEquals, sql.NullString{Valid: true, String: `darth@vader.com`}},
		{`t_users.age`, Equals, int64(40)},
		{`t_users.active`, Equals, true},
		{`t_users.email_verified`, SQLEquals, sql.NullBool{Valid: true, Bool: true}},
		{`t_users.no_of_licenses`, SQLEquals, sql.NullInt64{Valid: true, Int64: int64(0)}},
		{`t_users.created_at`, SQLEquals, testTime},
		{`t_users.last_payment_at`, SQLEquals, pq.NullTime{Valid: false, Time: testTime}},
	}
	recordCheck(data[0], tests, c)
}

func (s *UpdateExecTS) SetUpTest(c *C) {
	createTestTables()
	Register(`t_users`)
	_, err := Exec(sampleInsert)
	c.Assert(err, IsNil)

	s.query = &Query{
		queryType: UpdateQuery,
		query:     `UPDATE t_users SET email = $1, active = $2, email_verified = $3 WHERE id = $4`,
		args:      []interface{}{`luke@skywalker.com`, false, nil, `2u`},
	}
}

func (s *UpdateExecTS) TestUpdateTest(c *C) {
	data, err := Exec(s.query)
	c.Assert(err, IsNil)
	c.Assert(data, IsNil)

	data, _ = Exec(sampleSelect)
	var tests = []testEntry{
		{`t_users.id`, Equals, `2u`},
		{`t_users.email`, SQLEquals, sql.NullString{Valid: true, String: `luke@skywalker.com`}},
		{`t_users.age`, Equals, int64(40)},
		{`t_users.active`, Equals, false},
		{`t_users.email_verified`, SQLEquals, sql.NullBool{Valid: false, Bool: false}},
		{`t_users.no_of_licenses`, SQLEquals, sql.NullInt64{Valid: true, Int64: int64(0)}},
		{`t_users.created_at`, SQLEquals, testTime},
		{`t_users.last_payment_at`, SQLEquals, pq.NullTime{Valid: false, Time: testTime}},
	}
	recordCheck(data[0], tests, c)
}

func (s *DeleteExecTS) SetUpTest(c *C) {
	createTestTables()
	Register(`t_users`)
	_, err := Exec(sampleInsert)
	c.Assert(err, IsNil)

	s.query = &Query{
		queryType: DeleteQuery,
		query:     `DELETE FROM t_users WHERE id = $1`,
		args:      []interface{}{`2u`},
	}
}

func (s *DeleteExecTS) TestDeleteExec(c *C) {
	data, err := Exec(s.query)
	c.Assert(err, IsNil)
	c.Assert(data, IsNil)

	data, _ = Exec(sampleSelect)
	c.Assert(len(data), Equals, 0)
}

// TODO: add Truncate test, dry this up
