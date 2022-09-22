package model

import "math"

type Pulse struct {
	max float64
	td  float64
	tr  float64
	tf  float64
	pw  float64
	per float64
}

type Current struct {
	name  string
	val   float64
	pulse Pulse
}

// Create new current source for dc modeling
func NewCurrent(name string, val float64) *Current {
	_source := &Current{name: name, val: val}
	return _source
}

// Create new current source for tr modeling
func NewCurrentPulse(name string, val float64, max float64, td float64, tr float64, tf float64, pw float64, per float64) *Current {
	_source := &Current{name: name, val: val, pulse: Pulse{max, td, tr, tf, pw, per}}
	return _source
}

// Return valu of pulse current source in given time
func (c *Current) PulseValue(pulseTime float64) (float64, float64) {
	if pulseTime > c.pulse.td {
		pulseTime -= c.pulse.td
		pulseTime -= math.Floor(pulseTime/c.pulse.per) * c.pulse.per

		if pulseTime < c.pulse.tr {
			kr := (c.pulse.max - c.val) / (c.pulse.tr - 0)
			br := c.val - kr*0

			return kr*pulseTime + br, kr
		} else if pulseTime < c.pulse.tr+c.pulse.pw {
			return c.val, 0
		} else if pulseTime < c.pulse.tr+c.pulse.pw+c.pulse.tf {
			bf := (c.val*(c.pulse.tr+c.pulse.tf+c.pulse.pw) - c.pulse.max*(c.pulse.tr+c.pulse.pw)) / (-(c.pulse.tr + c.pulse.tf + c.pulse.pw) + (c.pulse.tr + c.pulse.pw))
			kf := (c.val - bf) / (c.pulse.tr + c.pulse.pw)

			return kf*pulseTime + bf, kf
		}
	}

	return c.val, 0
}
