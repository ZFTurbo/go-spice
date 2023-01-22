package model

type Current struct {
	name  string
	value float64
}

// Create new current source for dc modeling
func NewCurrent(name string, value float64) *Current {
	source := &Current{name: name, value: value}

	return source
}
