package simulation

import "github.com/erikpa1/TurtleIntelligenceBackend/tools"

type SimStepper struct {
	Now tools.Seconds
	End tools.Seconds
}

func (self *SimStepper) Step() {
	self.Now += 1
}
