package sqlparserhelper

/* sqlparser helper to parse SELECT / INSERT sql queries */

import (
	"github.com/xwb1989/sqlparser"
)

type Row struct {
	// column types are: string, int, float32
	Columns interface{}
}

const (
	CONDITION_OP_EQUALS = iota
	CONDITION_OP_GT
	CONDITION_OP_LT
	CONDITION_OP_GTE
	CONDITION_OP_LTE
)

const (
	COMMAND_TYPE_SELECT = iota
	COMMANT_TYPE_INSERT
	COMMAND_TYPE_UNSUPPORTED = -1
)

func IsSupportedQuery(query sqlparser.Statement) bool {

	// from SELECT and INSERT we support only:

	//  SELECT
	//  * / key,value
	//  FROM (single table name only)
	//  WHERE
	//  AND / OR / IN / NOT IN / NOT / ()
	//  = / < / > / <= / >=
	//  LIMIT n

	//  INSERT INTO
	//  single table name only
	//  (keys) values (values), and possible alternatives of keys-values

	// all else - return false ?

	// low priority function to get implemented

	return true
}

func GetCommandType(query sqlparser.Statement) int {

	// COMMAND_TYPE_SELECT or COMMAND_TYPE_INSERT or COMMAND_TYPE_UNSUPPORTED

	return COMMAND_TYPE_SELECT
}

func GetFieldList(query sqlparser.Statement) []string {

	// array of fields, aliases not supported

	return make([]string, 0)
}

func GetTableName(query sqlparser.Statement) string {

	// only single table name supported, aliases not supported

	return ""
}

func GetResults(query sqlparser.Statement, data []Row) []Row {
	resultData := make([]Row, 0)

	// go through "data" array and put to resultData rows that match sqlparser.Statement.Where conditions
	// example SQL WHERE query part: SELECT * FROM table1 WHERE a = 1 AND b = 2 AND c = 3 AND (q = 1 OR q = 2) OR n = 2 AND p = 2;
	// ...
	//

	return resultData
}

func GetLimit(query sqlparser.Statement) int {
	return 0
}
