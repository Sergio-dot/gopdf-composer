package models

import "strings"

// RuntimeContext holds the runtime data available for variable substitution
// and condition evaluation during PDF generation.
type RuntimeContext struct {
	Data map[string]any
}

// Get retrieves a top-level value from the context by key.
func (rc *RuntimeContext) Get(field string) (any, bool) {
	val, exists := rc.Data[field]
	return val, exists
}

// GetNested traverses a dot-separated path (e.g., "user.profile.email") into
// nested map[string]any values and returns the value if found.
func (rc *RuntimeContext) GetNested(path string) (any, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, false
	}

	val, exists := rc.Data[parts[0]]
	if !exists {
		return nil, false
	}

	for _, part := range parts[1:] {
		m, ok := val.(map[string]any)
		if !ok {
			return nil, false
		}
		val, ok = m[part]
		if !ok {
			return nil, false
		}
	}

	return val, true
}

// Set assigns a value to the context at the given top-level key.
func (rc *RuntimeContext) Set(key string, value any) {
	rc.Data[key] = value
}

// Delete removes a top-level key from the context.
func (rc *RuntimeContext) Delete(key string) {
	delete(rc.Data, key)
}
