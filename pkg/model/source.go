package model

type Pulse struct {
	min float64
	tr  float64
	tf  float64
	pw  float64
	dl  float64
}

type Current struct {
	name  string
	val   float64
	pulse Pulse
}

func NewCurrent(name string, val float64) *Current {
	_source := &Current{name: name, val: val}
	return _source
}

func NewCurrentPulse(name string, val float64, min float64, tr float64, tf float64, pw float64, dl float64) *Current {
	_source := &Current{name: name, val: val, pulse: Pulse{min, tr, tf, pw, dl}}
	return _source
}

func (c *Current) AddVal(val float64) {
	c.val += val
}
