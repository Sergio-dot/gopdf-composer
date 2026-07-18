package models

import "strings"

type RuntimeContext struct {
	Data map[string]any
}

func (rc *RuntimeContext) Get(field string) (any, bool) {
	val, exists := rc.Data[field]
	return val, exists
}

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

func (rc *RuntimeContext) Set(key string, value any) {
	rc.Data[key] = value
}

func (rc *RuntimeContext) Delete(key string) {
	delete(rc.Data, key)
}
