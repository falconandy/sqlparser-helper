package sqlparserhelper

import (
	"github.com/stretchr/testify/assert"
	"github.com/xwb1989/sqlparser"
	"testing"
)

type testData struct {
	Id int
	A  int
}

func TestNil(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE a = 1"
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
		{testData{Id: 1, A: 1}},
		{testData{Id: 1, A: 2}},
		{testData{Id: 3, A: 1}},
	})
	assert.Equal(t, 3, len(results))
}

func TestWhereWithEqualComparison(t *testing.T) {
	queryStr := "SELECT * FROM table1 WHERE a = 1"
	query, err := sqlparser.Parse(queryStr)
	assert.Nil(t, err)
	results := GetResults(query, []Row{
		{testData{Id: 1, A: 1}},
		{testData{Id: 2, A: 2}},
		{testData{Id: 3, A: 1}},
	})
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 1, results[0].Columns.(testData).Id)
	assert.Equal(t, 3, results[1].Columns.(testData).Id)
}
