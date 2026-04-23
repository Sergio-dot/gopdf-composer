package evaluator

import (
	"fmt"
	"reflect"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func Evaluate(condition *models.Condition, runtimeCtx *models.RuntimeContext) (bool, error) {
	if condition == nil {
		return true, nil
	}

	fieldValue, exists := runtimeCtx.Get(condition.Field)
	if !exists {
		return false, fmt.Errorf("field not found %s", condition.Field)
	}

	switch condition.Op {
	case "==":
		return reflect.DeepEqual(fieldValue, condition.Value), nil
	case "!=":
		return !reflect.DeepEqual(fieldValue, condition.Value), nil
		// TODO: extend with more operators (>, <, >=, <=, in, contains)
	}

	return false, fmt.Errorf("unexpected condition: %v on context: %v", condition, runtimeCtx)
}
