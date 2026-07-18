package evaluator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func Evaluate(condition *models.Condition, runtimeCtx *models.RuntimeContext) (bool, error) {
	if condition == nil {
		return true, nil
	}

	if condition.And != nil {
		for _, c := range condition.And {
			result, err := Evaluate(&c, runtimeCtx)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
	}

	if condition.Or != nil {
		for _, c := range condition.Or {
			result, err := Evaluate(&c, runtimeCtx)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
	}

	if condition.Not != nil {
		result, err := Evaluate(condition.Not, runtimeCtx)
		if err != nil {
			return false, err
		}
		return !result, nil
	}

	return evaluateLeaf(condition, runtimeCtx)
}

func evaluateLeaf(condition *models.Condition, runtimeCtx *models.RuntimeContext) (bool, error) {
	fieldValue, exists := runtimeCtx.Get(condition.Field)
	if !exists {
		return false, fmt.Errorf("field not found: %s", condition.Field)
	}

	switch condition.Op {
	case "==":
		return reflect.DeepEqual(fieldValue, condition.Value), nil
	case "!=":
		return !reflect.DeepEqual(fieldValue, condition.Value), nil
	case ">", "<", ">=", "<=":
		a, b, err := toFloat64Pair(fieldValue, condition.Value)
		if err != nil {
			return false, fmt.Errorf("%s requires numeric values: got %T and %T", condition.Op, fieldValue, condition.Value)
		}
		switch condition.Op {
		case ">":
			return a > b, nil
		case "<":
			return a < b, nil
		case ">=":
			return a >= b, nil
		case "<=":
			return a <= b, nil
		}
	case "in":
		return evalIn(fieldValue, condition.Value)
	case "contains":
		a, okA := fieldValue.(string)
		b, okB := condition.Value.(string)
		if !okA || !okB {
			return false, fmt.Errorf("contains requires string values")
		}
		return strings.Contains(a, b), nil
	}

	return false, fmt.Errorf("unknown operator: %s", condition.Op)
}

func toFloat64Pair(a, b any) (float64, float64, error) {
	fa, ok := toFloat64(a)
	if !ok {
		return 0, 0, fmt.Errorf("cannot convert %T to float64", a)
	}
	fb, ok := toFloat64(b)
	if !ok {
		return 0, 0, fmt.Errorf("cannot convert %T to float64", b)
	}
	return fa, fb, nil
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	}
	return 0, false
}

func evalIn(fieldValue, condValue any) (bool, error) {
	rv := reflect.ValueOf(condValue)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false, fmt.Errorf("in requires an array value, got %T", condValue)
	}

	for i := 0; i < rv.Len(); i++ {
		if reflect.DeepEqual(fieldValue, rv.Index(i).Interface()) {
			return true, nil
		}
	}
	return false, nil
}
