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

func TestNil(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results := GetResults(query, nil)
	assert.Equal(t, 0, len(results))
}

func TestWithoutWhere(t *testing.T) {
	queryStr := "SELECT * FROM table1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results := GetResults(query, []Row{
		{testData{Id: 1, I: 1}},
		{testData{Id: 1, I: 2}},
		{testData{Id: 3, I: 1}},
	})
	assert.Equal(t, 3, len(results))
}

func TestWhereWithEqualComparison(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results := GetResults(query, []Row{
		{testData{Id: 1, I: 1}},
		{testData{Id: 2, I: 2}},
		{testData{Id: 3, I: 1}},
	})
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 3, results[1].Columns.(testData).Id)
}

func TestWhereAnd(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1 and s != 'a'"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "a"}},
		{testData{Id: 2, I: 2, S: "b"}},
		{testData{Id: 3, I: 1, S: "c"}},
	})
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 3, results[0].Columns.(testData).Id)
}

func TestWhereOr(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE i = 1 or s != 'a'"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results := GetResults(query, []Row{
		{testData{Id: 1, I: 1, S: "a"}},
		{testData{Id: 2, I: 2, S: "b"}},
		{testData{Id: 3, I: 3, S: "a"}},
	})
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 2, results[1].Columns.(testData).Id)
}
