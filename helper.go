package sqlparserhelper

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
    //  AND / OR / IN / NOT IN / NOT
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
    return make([]Condition, 0)
}

func GetLimit(&(sqlparser.Statement)) int {
    return 0
}
