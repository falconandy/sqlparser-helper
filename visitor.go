package sqlparserhelper

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"reflect"
	"strconv"
	"strings"
)

type WhereVisitor struct {
	where sqlparser.Expr
}

type Predicate = func(item interface{}) (bool, error)

func (v *WhereVisitor) Visit() Predicate {
	return v.visitExpr(v.where)
}

func (v *WhereVisitor) visitExpr(expr sqlparser.Expr) Predicate {
	switch expr := expr.(type) {
	case *sqlparser.AndExpr:
		return v.visitAndExpr(expr)
	case *sqlparser.OrExpr:
		return v.visitOrExpr(expr)
	case *sqlparser.ComparisonExpr:
		return v.visitComparisonExpr(expr)
	case *sqlparser.ParenExpr:
		return v.visitParenExpr(expr)
	}

	return func(item interface{}) (bool, error) {
		return false, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (v *WhereVisitor) visitAndExpr(expr *sqlparser.AndExpr) Predicate {
	return func(item interface{}) (bool, error) {
		leftPredicate := v.visitExpr(expr.Left)
		left, err := leftPredicate(item)
		if err != nil {
			return false, err
		}

		if !left {
			return false, nil
		}

		rightPredicate := v.visitExpr(expr.Right)
		right, err := rightPredicate(item)
		if err != nil {
			return false, err
		}
		return right, nil
	}
}

func (v *WhereVisitor) visitOrExpr(expr *sqlparser.OrExpr) Predicate {
	return func(item interface{}) (bool, error) {
		leftPredicate := v.visitExpr(expr.Left)
		left, err := leftPredicate(item)
		if err != nil {
			return false, err
		}

		if left {
			return true, nil
		}

		rightPredicate := v.visitExpr(expr.Right)
		right, err := rightPredicate(item)
		if err != nil {
			return false, err
		}
		return right, nil
	}
}

func (v *WhereVisitor) visitComparisonExpr(expr *sqlparser.ComparisonExpr) Predicate {
	return func(item interface{}) (bool, error) {
		left, err := v.getValue(item, expr.Left)
		if err != nil {
			return false, err
		}

		right, err := v.getValue(item, expr.Right)
		if err != nil {
			return false, err
		}

		leftValue := reflect.ValueOf(left)
		rightValue := reflect.ValueOf(right)
		switch leftValue.Kind() {
		case reflect.String:
			switch rightValue.Kind() {
			case reflect.String:
				return v.compareStrings(leftValue.String(), rightValue.String(), expr.Operator)
			}
		case reflect.Int:
			switch rightValue.Kind() {
			case reflect.Int:
				return v.compareInts(int(leftValue.Int()), int(rightValue.Int()), expr.Operator)
			case reflect.Float64:
				return v.compareFloats(float64(leftValue.Int()), rightValue.Float(), expr.Operator)
			}
		case reflect.Float64:
			switch rightValue.Kind() {
			case reflect.Int:
				return v.compareFloats(leftValue.Float(), float64(rightValue.Int()), expr.Operator)
			case reflect.Float64:
				return v.compareFloats(leftValue.Float(), rightValue.Float(), expr.Operator)
			}
		}
		return false, fmt.Errorf("unsupported comparison: %v and %v", leftValue.Type(), rightValue.Type())
	}
}

func (v *WhereVisitor) visitParenExpr(expr *sqlparser.ParenExpr) Predicate {
	return v.visitExpr(expr.Expr)
}

func (v *WhereVisitor) compareStrings(left, right, operator string) (bool, error) {
	switch operator {
	case sqlparser.EqualStr:
		return left == right, nil
	case sqlparser.NotEqualStr:
		return left != right, nil
	case sqlparser.LessThanStr:
		return left < right, nil
	case sqlparser.LessEqualStr:
		return left <= right, nil
	case sqlparser.GreaterThanStr:
		return left > right, nil
	case sqlparser.GreaterEqualStr:
		return left >= right, nil
	}
	return false, fmt.Errorf("unsupported operator: %s", operator)
}

func (v *WhereVisitor) compareInts(left, right int, operator string) (bool, error) {
	switch operator {
	case sqlparser.EqualStr:
		return left == right, nil
	case sqlparser.NotEqualStr:
		return left != right, nil
	case sqlparser.LessThanStr:
		return left < right, nil
	case sqlparser.LessEqualStr:
		return left <= right, nil
	case sqlparser.GreaterThanStr:
		return left > right, nil
	case sqlparser.GreaterEqualStr:
		return left >= right, nil
	}
	return false, fmt.Errorf("unsupported operator: %s", operator)
}

func (v *WhereVisitor) compareFloats(left, right float64, operator string) (bool, error) {
	switch operator {
	case sqlparser.EqualStr:
		return left == right, nil
	case sqlparser.NotEqualStr:
		return left != right, nil
	case sqlparser.LessThanStr:
		return left < right, nil
	case sqlparser.LessEqualStr:
		return left <= right, nil
	case sqlparser.GreaterThanStr:
		return left > right, nil
	case sqlparser.GreaterEqualStr:
		return left >= right, nil
	}
	return false, fmt.Errorf("unsupported operator: %s", operator)
}

func (v *WhereVisitor) getValue(item interface{}, expr sqlparser.Expr) (interface{}, error) {
	valueExpr, ok := expr.(*sqlparser.SQLVal)
	if ok {
		switch valueExpr.Type {
		case sqlparser.StrVal:
			return string(valueExpr.Val), nil
		case sqlparser.IntVal:
			value, err := strconv.Atoi(string(valueExpr.Val))
			if err != nil {
				return nil, fmt.Errorf("can't parse int: %v", err)
			} else {
				return value, nil
			}
		case sqlparser.FloatVal:
			value, err := strconv.ParseFloat(string(valueExpr.Val), 64)
			if err != nil {
				return nil, fmt.Errorf("can't parse float: %v", err)
			} else {
				return value, nil
			}
		}
		return nil, fmt.Errorf("unsupported value type: %d", valueExpr.Type)
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
					return int(fieldValue.Int()), nil
				case reflect.Float32, reflect.Float64:
					return fieldValue.Float(), nil
				case reflect.String:
					return fieldValue.String(), nil
				}
			}
			return nil, fmt.Errorf("unsupported column type: %v", fieldValue.Type())
		}
	}

	return nil, fmt.Errorf("unsupported expression type: %T", expr)
}
