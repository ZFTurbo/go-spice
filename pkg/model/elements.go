package model

type Capaster struct {
	end interface{}
	val float64
}

func NewCapaster(end string, val float64) *Capaster {
	cap := &Capaster{end: end, val: val}
	return cap
}

type Inductance struct {
	end interface{}
	val float64
}

func NewInductance(end string, val float64) *Inductance {
	indc := &Inductance{end: end, val: val}
	return indc
}
