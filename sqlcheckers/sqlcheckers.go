// Package customcheckers implements a custom checker for Gocheck (gopkg.in/check.v1), which works on nullable sql type (sql.NulLString, pq.NullTime, etc)
package sqlcheckers

import (
	`database/sql`
	`fmt`
	`github.com/lib/pq`
	gocheck `gopkg.in/check.v1`
	`reflect`
	`time`
)

type sqlEqualsChecker struct {
	*gocheck.CheckerInfo
}

const (
	unsupportedType = `unsupported type %v, type must be one of {'sql.NullString', 'sql.NullBool'. 'sql.NullInt64'. 'pq.NullTime'}`
	typeMismatch    = `type mismatched: obtained type %v, comparing to type %v`
	invalidMsg      = `Valid not equal, obtained.Valid = %v, expected.Valid = %v`
)

// SQLEquals is a checker that helps checking equality of sql.NullString, sql.NullBool, sql.NullInt64 and pq.NullTime
var SQLEquals = &sqlEqualsChecker{
	&gocheck.CheckerInfo{
		Name:   `SQLEquals`,
		Params: []string{`obtained`, `expected`},
	},
}

func (checker *sqlEqualsChecker) Check(params []interface{}, names []string) (result bool, err string) {
	obtained, expected := params[0], params[1]

	if reflect.TypeOf(obtained) != reflect.TypeOf(expected) {
		return false, fmt.Sprintf(typeMismatch, reflect.TypeOf(obtained), reflect.TypeOf(expected))
	}

	switch obtained.(type) {
	case sql.NullString:
		return nullStringEqual(obtained.(sql.NullString), expected.(sql.NullString))
	case sql.NullBool:
		return nullBoolEqual(obtained.(sql.NullBool), expected.(sql.NullBool))
	case sql.NullInt64:
		return nullInt64Equal(obtained.(sql.NullInt64), expected.(sql.NullInt64))
	case pq.NullTime:
		return nullTimeEqual(obtained.(pq.NullTime), expected.(pq.NullTime))
	case time.Time:
		return timeEqual(obtained.(time.Time), expected.(time.Time))
	}

	return false, fmt.Sprintf(unsupportedType, reflect.TypeOf(obtained))
}

func compareValue(obtained, expected interface{}) (result bool, err string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			err = fmt.Sprint(v)
		}
	}()
	return obtained == expected, ``
}

type comparator func(interface{}, interface{}) (bool, string)

// return (done, result, err), indicating if comparison can be stopped base on .Valid & the actual results
func compareValid(obtained, expected bool) (bool, bool, string) {
	if obtained != expected {
		return true, false, fmt.Sprintf(invalidMsg, obtained, expected)
	}

	if obtained == false {
		return true, true, ``
	}

	return false, false, `` // continue to compare value
}

func nullStringEqual(obtained, expected sql.NullString) (result bool, err string) {
	if done, res, err := compareValid(obtained.Valid, expected.Valid); done {
		return res, err
	}
	return compareValue(obtained.String, expected.String)
}

func nullBoolEqual(obtained, expected sql.NullBool) (result bool, err string) {
	if done, res, err := compareValid(obtained.Valid, expected.Valid); done {
		return res, err
	}
	return compareValue(obtained.Bool, expected.Bool)
}

func nullInt64Equal(obtained, expected sql.NullInt64) (result bool, err string) {
	if done, res, err := compareValid(obtained.Valid, expected.Valid); done {
		return res, err
	}
	return compareValue(obtained.Int64, expected.Int64)
}

// for time, just check if when convert to UTC, the no of seconds is the same
// otherwise, we need to check for the location in time.Time as well
// this location changes depending the timezone of the maching running the test
// so it's not nice to use
func timeEqual(obtained, expected time.Time) (result bool, err string) {
	return compareValue(obtained.UTC().Unix(), expected.UTC().Unix())
}

func nullTimeEqual(obtained, expected pq.NullTime) (result bool, err string) {
	if done, res, err := compareValid(obtained.Valid, expected.Valid); done {
		return res, err
	}
	return compareValue(obtained.Time.UTC().Unix(), expected.Time.UTC().Unix())
}
