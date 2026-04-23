package models

type RuntimeContext struct {
	Data map[string]any
}

func (rc *RuntimeContext) Get(field string) (any, bool) {
	val, exists := rc.Data[field]
	return val, exists
}
