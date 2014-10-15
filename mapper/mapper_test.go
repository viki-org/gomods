package mapper

import (
	. `gopkg.in/check.v1`
	`testing`
)

func Test(t *testing.T) {
	TestingT(t)
}

type SchemaTS struct{}

type SchemaRegisterTS struct {
	query Query
}

func init() {
	Suite(&SchemaTS{})
	Suite(&SchemaRegisterTS{})
}

func (s *SchemaRegisterTS) SetUpTest(c *C) {
	createTestTables()
}

func (s *SchemaRegisterTS) TestRegister(c *C) {
	Register(`t_users`)

	c.Assert(len(columns), Equals, 8)
	var tests = []testEntry{
		{`t_users.id`, Equals, stringType},
		{`t_users.email`, Equals, nullStringType},
		{`t_users.age`, Equals, int64Type},
		{`t_users.active`, Equals, boolType},
		{`t_users.email_verified`, Equals, nullBoolType},
		{`t_users.no_of_licenses`, Equals, nullInt64Type},
		{`t_users.last_payment_at`, Equals, nullTimeType},
		{`t_users.created_at`, Equals, timeType},
	}
	tableCheck(c, tests, func(target interface{}) interface{} {
		return columns[target.(string)]
	})
}

var toTypeTests = []struct {
	dataType, nullable string
	out                interface{}
}{
	{`character varying`, `NO`, stringType},
	{`character varying`, `YES`, nullStringType},
	{`text`, `NO`, stringType},
	{`text`, `YES`, nullStringType},
	{`integer`, `NO`, int64Type},
	{`integer`, `YES`, nullInt64Type},
	{`boolean`, `NO`, boolType},
	{`boolean`, `YES`, nullBoolType},
	{`timestamp with time zone`, `NO`, timeType},
	{`timestamp with time zone`, `YES`, nullTimeType},
	{`timestamp without time zone`, `NO`, timeType},
	{`timestamp without time zone`, `YES`, nullTimeType},
	{`random`, `YES`, `invalid sql data type, got "random", expected one of ("character varying", "text", "integer", "boolean", "timestamp with time zone", "timestamp without time zone")`},
	{`character varying`, `yes`, `invalid value for nullable, got "yes", expected one of ("YES", "NO")`},
}

func (s *SchemaTS) TestToType(c *C) {
	for _, test := range toTypeTests {
		golangType, err := toType(test.dataType, test.nullable)
		if err != nil {
			c.Assert(err.Error(), Equals, test.out.(string))
		} else {
			c.Assert(golangType, Equals, test.out.(int))
		}
	}
}
