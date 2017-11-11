package sqlparserhelper

import (
	"github.com/stretchr/testify/assert"
	"github.com/xwb1989/sqlparser"
	"testing"
)

type testData struct {
	Id int
	I  int
	S  string
	F  float32
	D  float64
}

func TestGetResultsOfNil(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(results))
}

func TestGetResultsOfInsert(t *testing.T) {
	queryStr := "INSERT INTO table1 values(1)"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	_, err = GetResults(query, nil)
	assert.EqualError(t, err, "unsupported query type: *sqlparser.Insert")
}

func TestGetResultsWithoutWhere(t *testing.T) {
	queryStr := "SELECT * FROM table1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1}},
		{testData{Id: 1, I: 2}},
		{testData{Id: 3, I: 1}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(results))
}

func TestGetResultsWithUnsupportedExpression(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE EXISTS(select 1 from table2)"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	_, err = GetResults(query, []Row{
		{testData{Id: 1, I: 1}},
	})
	assert.EqualError(t, err, "unsupported expression type: *sqlparser.ExistsExpr")
}

func TestGetResultsWhereInt(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1}},
		{testData{Id: 2, I: 2}},
		{testData{Id: 3, I: 1}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 3, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereString(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE s != 'qwerty'"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, S: "qwerty"}},
		{testData{Id: 2, S: "1234"}},
		{testData{Id: 3, S: "mmm"}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 2, results[0].Columns.(testData).Id)
	assert.Equal(t, 3, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereFloat(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE f < 2.6"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, F: 1.0}},
		{testData{Id: 2, F: 2.5}},
		{testData{Id: 3, F: 3.3}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 2, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereUnsupportedValue(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i < 0x1234"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	_, err = GetResults(query, []Row{
		{testData{Id: 1, I: 2}},
	})
	assert.EqualError(t, err, "unsupported value type: 3")
}

func TestGetResultsWhereAnd(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1 and s != 'a'"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "a"}},
		{testData{Id: 2, I: 2, S: "b"}},
		{testData{Id: 3, I: 1, S: "c"}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 3, results[0].Columns.(testData).Id)
}

func TestGetResultsWhereOr(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE f >= 2.5 or s = 'a'"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, F: 2.5, S: "a"}},
		{testData{Id: 2, F: 2, S: "a"}},
		{testData{Id: 3, F: 1.8, S: "b"}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 2, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereAndOr(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i != 1 and s <= 'd' or f > 3"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "a", F: 4}},
		{testData{Id: 2, I: 2, S: "f", F: 3.6}},
		{testData{Id: 3, I: 3, S: "z", F: 3}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 2, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereOrAnd(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i <= 1.5 or s <= 'd' and f > 3"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "a", F: 4}},
		{testData{Id: 2, I: 2, S: "b", F: 3.6}},
		{testData{Id: 3, I: 3, S: "z", F: 3}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 2, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereParentheses(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE (i <= 1.5 or s <= 'd') and f > 3"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "a", F: 4}},
		{testData{Id: 2, I: 2, S: "b", F: 3.6}},
		{testData{Id: 3, I: 3, S: "z", F: 3}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 2, results[1].Columns.(testData).Id)
}

func TestGetResultsWhereComplex(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE I = 1 AND S = 'qwe' AND F = 3.5 AND (d > 2 OR d < 1) OR i = 2 AND f >= 2.6"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results, err := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "qwe", F: 3.5, D: 7.3}},
		{testData{Id: 2, I: 2, S: "qwe", F: 2.5, D: -1}},
		{testData{Id: 3, I: 2, S: "z", F: 2.7, D: 3}},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 3, results[1].Columns.(testData).Id)
}
