package sqlparserhelper

/* sqlparser helper to parse SELECT / INSERT sql queries */

import (
    "github.com/xwb1989/sqlparser"
)

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
)

type Condition struct {
    Key string
    Operator int
    Value interface{}
}

func IsSupportedQuery(&(sqlparser.Statement)) bool {

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

    return true
}

func GetCommandType(&(sqlparser.Statement)) int {

    // COMMAND_TYPE_SELECT or COMMAND_TYPE_INSERT or -1

    return COMMAND_TYPE_SELECT
}

func GetFieldList(&(sqlparser.Statement)) []string {

    // array of fields, aliases not supported

    return make([]string, 0)
}

func GetTableName(&(sqlparser.Statement)) string {

    // only single table name supported, aliases not supported

    return ""
}

func GetConditions(&(sqlparser.Statement)) []Condition {

    // example transformation of stmt.Where into an array of AND statements
    // argument query: a = 1 AND b = 2 AND c = 3 AND (q = 1 OR q = 2) OR n = 2 AND p = 2
    // result[0]: a = 1 AND b = 2 AND c = 3 AND q = 1
    // result[1]: a = 1 AND b = 2 AND c = 3 AND q = 2
    // result[2]: n = 2 AND p = 2

    return make([]Condition, 0)
}

func GetLimit(&(sqlparser.Statement)) int {
    return 0
}
