package evaluator

import (
	"testing"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func TestEvaluateNilCondition(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	result, err := Evaluate(nil, rc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result {
		t.Error("nil condition should evaluate to true")
	}
}

func TestEvaluateLeaf(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"name":   "Sergio",
		"age":    float64(30),
		"active": true,
		"score":  float64(85),
	}}
	ctxWithSlice := &models.RuntimeContext{Data: map[string]any{
		"name":    "Sergio",
		"country": "ES",
		"tags":    []any{"go", "pdf", "composer"},
	}}

	tests := []struct {
		name      string
		rc        *models.RuntimeContext
		condition *models.Condition
		want      bool
		wantErr   bool
	}{
		// == / != operators
		{name: "eq string match", rc: rc, condition: &models.Condition{Field: "name", Op: "==", Value: "Sergio"}, want: true},
		{name: "eq string no match", rc: rc, condition: &models.Condition{Field: "name", Op: "==", Value: "Maria"}, want: false},
		{name: "eq number match", rc: rc, condition: &models.Condition{Field: "age", Op: "==", Value: float64(30)}, want: true},
		{name: "eq number no match", rc: rc, condition: &models.Condition{Field: "age", Op: "==", Value: float64(25)}, want: false},
		{name: "eq bool match", rc: rc, condition: &models.Condition{Field: "active", Op: "==", Value: true}, want: true},
		{name: "ne string match", rc: rc, condition: &models.Condition{Field: "name", Op: "!=", Value: "Maria"}, want: true},
		{name: "ne string no match", rc: rc, condition: &models.Condition{Field: "name", Op: "!=", Value: "Sergio"}, want: false},

		// Numeric comparison
		{name: "gt number true", rc: rc, condition: &models.Condition{Field: "score", Op: ">", Value: float64(80)}, want: true},
		{name: "gt number false", rc: rc, condition: &models.Condition{Field: "score", Op: ">", Value: float64(90)}, want: false},
		{name: "lt number true", rc: rc, condition: &models.Condition{Field: "score", Op: "<", Value: float64(90)}, want: true},
		{name: "lt number false", rc: rc, condition: &models.Condition{Field: "score", Op: "<", Value: float64(80)}, want: false},
		{name: "gte number true", rc: rc, condition: &models.Condition{Field: "score", Op: ">=", Value: float64(85)}, want: true},
		{name: "gte number false", rc: rc, condition: &models.Condition{Field: "score", Op: ">=", Value: float64(86)}, want: false},
		{name: "lte number true", rc: rc, condition: &models.Condition{Field: "score", Op: "<=", Value: float64(85)}, want: true},
		{name: "lte number false", rc: rc, condition: &models.Condition{Field: "score", Op: "<=", Value: float64(84)}, want: false},

		// Numeric comparison with int context value
		{name: "gt with int context", rc: &models.RuntimeContext{Data: map[string]any{"count": 10}},
			condition: &models.Condition{Field: "count", Op: ">", Value: float64(5)}, want: true},
		{name: "lt int true", rc: &models.RuntimeContext{Data: map[string]any{"count": 10}},
			condition: &models.Condition{Field: "count", Op: "<", Value: float64(15)}, want: true},

		// Numeric comparison with non-numeric values
		{name: "gt non-numeric error", rc: rc,
			condition: &models.Condition{Field: "name", Op: ">", Value: "x"}, wantErr: true},

		// in operator
		{name: "in true", rc: ctxWithSlice,
			condition: &models.Condition{Field: "country", Op: "in", Value: []any{"ES", "PT", "FR"}}, want: true},
		{name: "in false", rc: ctxWithSlice,
			condition: &models.Condition{Field: "country", Op: "in", Value: []any{"US", "UK"}}, want: false},
		{name: "in with int", rc: &models.RuntimeContext{Data: map[string]any{"level": float64(3)}},
			condition: &models.Condition{Field: "level", Op: "in", Value: []any{float64(1), float64(2), float64(3)}}, want: true},
		{name: "in with non-array error", rc: rc,
			condition: &models.Condition{Field: "name", Op: "in", Value: "not-an-array"}, wantErr: true},

		// contains operator
		{name: "contains true", rc: rc,
			condition: &models.Condition{Field: "name", Op: "contains", Value: "erg"}, want: true},
		{name: "contains false", rc: rc,
			condition: &models.Condition{Field: "name", Op: "contains", Value: "xyz"}, want: false},
		{name: "contains non-string error", rc: rc,
			condition: &models.Condition{Field: "age", Op: "contains", Value: "x"}, wantErr: true},

		// Edge cases
		{name: "missing field", rc: rc,
			condition: &models.Condition{Field: "missing", Op: "==", Value: "x"}, wantErr: true},
		{name: "unknown operator", rc: rc,
			condition: &models.Condition{Field: "name", Op: "!!", Value: "x"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.condition, tt.rc)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}

func TestEvaluateCompoundAnd(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"age": float64(25), "country": "ES", "active": true,
	}}

	// All true → true
	cond := &models.Condition{
		And: []models.Condition{
			{Field: "age", Op: ">=", Value: float64(18)},
			{Field: "country", Op: "in", Value: []any{"ES", "PT"}},
			{Field: "active", Op: "==", Value: true},
		},
	}
	result, err := Evaluate(cond, rc)
	if err != nil || !result {
		t.Errorf("and-all-true: got %v, err %v, want true", result, err)
	}

	// One false → false
	cond = &models.Condition{
		And: []models.Condition{
			{Field: "age", Op: ">=", Value: float64(18)},
			{Field: "country", Op: "==", Value: "US"},
			{Field: "active", Op: "==", Value: true},
		},
	}
	result, err = Evaluate(cond, rc)
	if err != nil || result {
		t.Errorf("and-one-false: got %v, err %v, want false", result, err)
	}

	// Empty and → true
	cond = &models.Condition{And: []models.Condition{}}
	result, err = Evaluate(cond, rc)
	if err != nil || !result {
		t.Errorf("and-empty: got %v, err %v, want true", result, err)
	}
}

func TestEvaluateCompoundOr(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"country": "ES",
	}}

	// One true → true
	cond := &models.Condition{
		Or: []models.Condition{
			{Field: "country", Op: "==", Value: "US"},
			{Field: "country", Op: "==", Value: "ES"},
			{Field: "country", Op: "==", Value: "FR"},
		},
	}
	result, err := Evaluate(cond, rc)
	if err != nil || !result {
		t.Errorf("or-one-true: got %v, err %v, want true", result, err)
	}

	// All false → false
	cond = &models.Condition{
		Or: []models.Condition{
			{Field: "country", Op: "==", Value: "US"},
			{Field: "country", Op: "==", Value: "UK"},
		},
	}
	result, err = Evaluate(cond, rc)
	if err != nil || result {
		t.Errorf("or-all-false: got %v, err %v, want false", result, err)
	}

	// Empty or → false
	cond = &models.Condition{Or: []models.Condition{}}
	result, err = Evaluate(cond, rc)
	if err != nil || result {
		t.Errorf("or-empty: got %v, err %v, want false", result, err)
	}
}

func TestEvaluateCompoundNot(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"active": true,
	}}

	cond := &models.Condition{
		Not: &models.Condition{Field: "active", Op: "==", Value: true},
	}
	result, err := Evaluate(cond, rc)
	if err != nil || result {
		t.Errorf("not-true: got %v, err %v, want false", result, err)
	}

	cond = &models.Condition{
		Not: &models.Condition{Field: "active", Op: "==", Value: false},
	}
	result, err = Evaluate(cond, rc)
	if err != nil || !result {
		t.Errorf("not-false: got %v, err %v, want true", result, err)
	}
}

func TestEvaluateNestedCompound(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"role": "editor", "approved": false, "country": "ES",
	}}

	// Admin or (editor and approved) — editor is false → false
	cond := &models.Condition{
		Or: []models.Condition{
			{Field: "role", Op: "==", Value: "admin"},
			{
				And: []models.Condition{
					{Field: "role", Op: "==", Value: "editor"},
					{Field: "approved", Op: "==", Value: true},
				},
			},
		},
	}
	result, err := Evaluate(cond, rc)
	if err != nil || result {
		t.Errorf("nested-or-and-false: got %v, err %v, want false", result, err)
	}

	// Modify: editor is approved → true
	rc.Set("approved", true)
	result, err = Evaluate(cond, rc)
	if err != nil || !result {
		t.Errorf("nested-or-and-true: got %v, err %v, want true", result, err)
	}
}

func TestEvaluateBackwardCompatible(t *testing.T) {
	// Old-style conditions (Field/Op/Value only, no And/Or/Not) still work
	rc := &models.RuntimeContext{Data: map[string]any{
		"name": "Sergio", "age": float64(30),
	}}

	cond := &models.Condition{Field: "name", Op: "==", Value: "Sergio"}
	result, err := Evaluate(cond, rc)
	if err != nil || !result {
		t.Errorf("backward-compat: got %v, err %v, want true", result, err)
	}
}
