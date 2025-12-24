package simulation

import "turtle/tools"

type SimStepper struct {
	Now tools.Seconds
	End tools.Seconds
}

func (self *SimStepper) Step() {
	self.Now += 1
}
