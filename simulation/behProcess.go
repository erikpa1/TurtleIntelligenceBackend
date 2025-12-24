package simulation

import (
	"fmt"

	"turtle/lg"
	"turtle/simulation/simInternal"
	"turtle/simulation/stats"
	"turtle/tools"
)

type ProcessStates int8

const (
	PROC_STAT_IDLE    = 0
	PROC_STAT_BLOCKED = 1
	PROC_STAT_WORKING = 2
)

type ProcessBehaviour struct {
	World  *SimWorld
	Entity *SimEntity

	ActiveActor *SimActor

	ProcessTime   string
	ProcessFinish tools.Seconds

	ActiveState ProcessStates

	Statistics *stats.ProcessStats
}

func NewProcessBehaviour() ISimBehaviour {
	return &ProcessBehaviour{
		ActiveState: PROC_STAT_IDLE,
		Statistics:  stats.NewProcessStats(),
	}
}

func (self *ProcessBehaviour) __Interface() {
	procesB := ProcessBehaviour{}
	var _ ISimBehaviour = &procesB
	var _ ActorTakerBehaviour = &procesB
	var _ AIProvider = &procesB

}

// IBehaviour implementation
func (self *ProcessBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *ProcessBehaviour) Init1() {

}

func (self *ProcessBehaviour) Init2() {

}

func (self *ProcessBehaviour) Step() {

	now := self.World.Stepper.Now

	if self.ActiveState == PROC_STAT_IDLE {
		//DO nothing
	} else if self.ActiveState == PROC_STAT_BLOCKED {
		self._TryToPassEntityNext()
	} else if self.ActiveState == PROC_STAT_WORKING {
		if now >= self.ProcessFinish {
			self._FinishManufacturing()
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

func (self *ProcessBehaviour) GetEntity() *SimEntity {
	return self.Entity
}

func (self *ProcessBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = entity
	self.ProcessTime = entity.TypeData.GetString("processTime", "00:10")

}

// Taker behaviour
func (self *ProcessBehaviour) TakeActor(actor *SimActor) bool {
	canTake := self.CanTakeActor(actor)

	if canTake {
		self.ActiveActor = actor
		actor.UpdatePosition(self.Entity.Position)
		self._StartManufacturing()
	}

	return canTake
}

func (self *ProcessBehaviour) CanTakeActor(actor *SimActor) bool {
	return self.ActiveActor == nil
}

func (self *ProcessBehaviour) _StartManufacturing() {
	finishTime := tools.Seconds(tools.AnyExpr_CompileSeconds(self.ProcessTime, 10))
	self.ProcessFinish = finishTime + self.World.Stepper.Now
	self.ChangeState(PROC_STAT_WORKING)

	self.World.CreateUpcomingEvent(simInternal.SimUpcomingEvent{
		Id:     self.Entity.RuntimeId,
		Type:   simInternal.UPC_EVENT_FINISH,
		Second: self.ProcessFinish,
	})

	lg.LogOk("Started manufacturing")

}

func (self *ProcessBehaviour) _FinishManufacturing() {
	self.ProcessFinish = tools.Seconds(tools.MaxInt64())

	lg.LogOk("Finished manufacturing")

	self._TryToPassEntityNext()

}

func (self *ProcessBehaviour) ChangeState(newState ProcessStates) {
	self.ActiveState = newState
}

func (self *ProcessBehaviour) _TryToPassEntityNext() {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {

			taker, isITaker := connection.(ActorTakerBehaviour)

			if isITaker {
				taken := taker.TakeActor(self.ActiveActor)

				if taken {
					self.ActiveActor = nil
					self.ChangeState(PROC_STAT_IDLE)
					return
				}
			}
		}
	}

	self.ChangeState(PROC_STAT_BLOCKED)

}
func (self *ProcessBehaviour) GetAIDescription() string {
	return fmt.Sprintf(`
Type: Work center
Name: %s
ProcessTime: %s
`,
		self.Entity.Name,
		self.Entity.TypeData.GetString("processTime", "00:10"),
	)
}
