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

func TestEvaluate(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"name":    "Sergio",
		"age":     float64(30),
		"active":  true,
		"balance": 100.50,
	}}

	tests := []struct {
		name      string
		condition *models.Condition
		want      bool
		wantErr   bool
	}{
		{
			name:      "eq string match",
			condition: &models.Condition{Field: "name", Op: "==", Value: "Sergio"},
			want:      true,
		},
		{
			name:      "eq string no match",
			condition: &models.Condition{Field: "name", Op: "==", Value: "Maria"},
			want:      false,
		},
		{
			name:      "eq number match",
			condition: &models.Condition{Field: "age", Op: "==", Value: float64(30)},
			want:      true,
		},
		{
			name:      "eq number no match",
			condition: &models.Condition{Field: "age", Op: "==", Value: float64(25)},
			want:      false,
		},
		{
			name:      "eq bool match",
			condition: &models.Condition{Field: "active", Op: "==", Value: true},
			want:      true,
		},
		{
			name:      "ne string match",
			condition: &models.Condition{Field: "name", Op: "!=", Value: "Maria"},
			want:      true,
		},
		{
			name:      "ne string no match",
			condition: &models.Condition{Field: "name", Op: "!=", Value: "Sergio"},
			want:      false,
		},
		{
			name:      "missing field",
			condition: &models.Condition{Field: "missing", Op: "==", Value: "x"},
			wantErr:   true,
		},
		{
			name:      "unknown operator",
			condition: &models.Condition{Field: "name", Op: ">", Value: "x"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.condition, rc)
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
