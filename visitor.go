package sqlparserhelper

import (
	"github.com/xwb1989/sqlparser"
	"reflect"
	"strconv"
	"strings"
)

type WhereVisitor struct {
	where sqlparser.Expr
}

type Predicate = func(item interface{}) bool

func (v *WhereVisitor) Visit() Predicate {
	return v.visitExpr(v.where)
}

func (v *WhereVisitor) visitExpr(expr sqlparser.Expr) Predicate {
	comparisonExpr, ok := expr.(*sqlparser.ComparisonExpr)
	if ok {
		return v.visitCompareExpr(comparisonExpr)
	}

	return func(item interface{}) bool {
		return false
	}
}

func (v *WhereVisitor) visitCompareExpr(expr *sqlparser.ComparisonExpr) Predicate {
	return func(item interface{}) bool {
		if left, ok := v.getValue(item, expr.Left); ok {
			if right, ok := v.getValue(item, expr.Right); ok {
				leftValue := reflect.ValueOf(left)
				rightValue := reflect.ValueOf(right)
				switch leftValue.Kind() {
				case reflect.String:
					switch rightValue.Kind() {
					case reflect.String:
						result, ok := v.compareStrings(leftValue.String(), rightValue.String(), expr.Operator)
						if ok {
							return result
						}
					}
				case reflect.Int:
					switch rightValue.Kind() {
					case reflect.Int:
						result, ok := v.compareInts(int(leftValue.Int()), int(rightValue.Int()), expr.Operator)
						if ok {
							return result
						}
					case reflect.Float64:
						result, ok := v.compareFloats(float64(leftValue.Int()), rightValue.Float(), expr.Operator)
						if ok {
							return result
						}
					}
				case reflect.Float64:
					switch rightValue.Kind() {
					case reflect.Int:
						result, ok := v.compareFloats(leftValue.Float(), float64(rightValue.Int()), expr.Operator)
						if ok {
							return result
						}
					case reflect.Float64:
						result, ok := v.compareFloats(leftValue.Float(), rightValue.Float(), expr.Operator)
						if ok {
							return result
						}
					}
				}
			}
		}
		return false
	}
}

func (v *WhereVisitor) compareStrings(left, right, operator string) (bool, bool) {
	switch operator {
	case sqlparser.EqualStr:
		return left == right, true
	case sqlparser.NotEqualStr:
		return left != right, true
	case sqlparser.LessThanStr:
		return left < right, true
	case sqlparser.LessEqualStr:
		return left <= right, true
	case sqlparser.GreaterThanStr:
		return left > right, true
	case sqlparser.GreaterEqualStr:
		return left >= right, true
	}
	return false, false
}

func (v *WhereVisitor) compareInts(left, right int, operator string) (bool, bool) {
	switch operator {
	case sqlparser.EqualStr:
		return left == right, true
	case sqlparser.NotEqualStr:
		return left != right, true
	case sqlparser.LessThanStr:
		return left < right, true
	case sqlparser.LessEqualStr:
		return left <= right, true
	case sqlparser.GreaterThanStr:
		return left > right, true
	case sqlparser.GreaterEqualStr:
		return left >= right, true
	}
	return false, false
}

func (v *WhereVisitor) compareFloats(left, right float64, operator string) (bool, bool) {
	switch operator {
	case sqlparser.EqualStr:
		return left == right, true
	case sqlparser.NotEqualStr:
		return left != right, true
	case sqlparser.LessThanStr:
		return left < right, true
	case sqlparser.LessEqualStr:
		return left <= right, true
	case sqlparser.GreaterThanStr:
		return left > right, true
	case sqlparser.GreaterEqualStr:
		return left >= right, true
	}
	return false, false
}

func (v *WhereVisitor) getValue(item interface{}, expr sqlparser.Expr) (interface{}, bool) {
	valueExpr, ok := expr.(*sqlparser.SQLVal)
	if ok {
		switch valueExpr.Type {
		case sqlparser.StrVal:
			return string(valueExpr.Val), true
		case sqlparser.IntVal:
			value, err := strconv.Atoi(string(valueExpr.Val))
			if err == nil {
				return value, true
			}
		case sqlparser.FloatVal:
			value, err := strconv.ParseFloat(string(valueExpr.Val), 64)
			if err == nil {
				return value, true
			}
		}
		return nil, false
	}

	columnExpr, ok := expr.(*sqlparser.ColName)
	if ok {
		itemValue := reflect.ValueOf(item)
		if itemValue.Kind() == reflect.Struct {
			columnName := strings.ToLower(columnExpr.Name.String())
			fieldValue := itemValue.FieldByNameFunc(func(field string) bool {
				return strings.ToLower(field) == columnName
			})
			if fieldValue.IsValid() {
				switch fieldValue.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
					return int(fieldValue.Int()), true
				case reflect.Float32, reflect.Float64:
					return fieldValue.Float(), true
				case reflect.String:
					return fieldValue.String(), true
				}
			}
			return nil, false
		}
	}

	return nil, false
}
