package simulation2

import (
	"turtle/core/lgr"
	"turtle/simulation/simInternal"
	"turtle/simulation/stats"
	"turtle/tools"
)

type ProcessStates int8

const (
	PROC_STAT_IDLE    ProcessStates = 0
	PROC_STAT_BLOCKED ProcessStates = 1
	PROC_STAT_WORKING ProcessStates = 2
)

type BehProcess struct {
	World       *SimWorld
	Entity      *SimEntity
	ActiveActor *SimActor

	ProcessTime   string
	ProcessFinish tools.Seconds

	ActiveState ProcessStates

	Statistics *stats.ProcessStats
}

func GetBehProcess(entity *SimEntity) *BehProcess {
	return entity.Impl.(*BehProcess)
}

func (self *BehProcess) TakeActor(actor *SimActor) bool {
	canTake := self.CanTakeActor(actor)

	if canTake {
		self.ActiveActor = actor
		actor.UpdatePosition(self.Entity.Position)
		self._StartManufacturing()
	}

	return canTake
}

func (self *BehProcess) CanTakeActor(actor *SimActor) bool {
	return self.ActiveActor == nil
}

func (self *BehProcess) Step() {

	now := self.World.Stepper.Now

	if self.ActiveActor != nil {
		if self.ActiveState == PROC_STAT_IDLE {
			//DO nothing
		} else if self.ActiveState == PROC_STAT_BLOCKED {
			self._TryToPassEntityNext()
		} else if self.ActiveState == PROC_STAT_WORKING {
			if now >= self.ProcessFinish {
				self._FinishManufacturing()
			}
		}
	}

	if self.ActiveState == PROC_STAT_IDLE {
		self.Statistics.IdleTime += 1
	} else if self.ActiveState == PROC_STAT_BLOCKED {
		self.Statistics.BlockedTime += 1
	} else if self.ActiveState == PROC_STAT_WORKING {
		self.Statistics.ProcessTime += 1
	}
}

func (self *BehProcess) _StartManufacturing() {
	finishTime := tools.Seconds(tools.AnyExpr_CompileSeconds(self.ProcessTime, 10))
	self.ProcessFinish = finishTime + self.World.Stepper.Now
	self.ChangeState(PROC_STAT_WORKING)

	self.World.CreateUpcomingEvent(simInternal.SimUpcomingEvent{
		Id:     self.Entity.RuntimeId,
		Type:   simInternal.UPC_EVENT_FINISH,
		Second: self.ProcessFinish,
	})

	lgr.Ok("Started manufacturing")

}

func (self *BehProcess) _FinishManufacturing() {
	self.ProcessFinish = tools.Seconds(tools.MaxInt64())

	lgr.Ok("Finished manufacturing")

	self._TryToPassEntityNext()

}

func (self *BehProcess) ChangeState(newState ProcessStates) {
	self.ActiveState = newState
}

func (self *BehProcess) _TryToPassEntityNext() {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {
			taken := connection.TakeActor(self.ActiveActor)

			if taken {
				self.ActiveActor = nil
				self.ChangeState(PROC_STAT_IDLE)
				return
			}
		}
	}

	self.ChangeState(PROC_STAT_BLOCKED)

}
