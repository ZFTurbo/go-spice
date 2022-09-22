package model

type Analysis struct {
	name     string
	fullTime float64
	timeStep float64
}

func NewAnalysis(name string, fullTime float64, timeStep float64) *Analysis {
	analysis := &Analysis{name: name, fullTime: fullTime, timeStep: timeStep}
	return analysis
}

func (a *Analysis) SetName(val string) {
	a.name = val
}

func (a *Analysis) SetFullTime(val float64){
	a.fullTime = val
}

func (a *Analysis) SetTimeStep(val float64){
	a.timeStep = val
}