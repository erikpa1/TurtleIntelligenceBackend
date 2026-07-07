package behSpawn

import (
	"turtle/core/lgr"
	"turtle/simulation/simInternal"
	"turtle/simulation2/entities"
	"turtle/simulation2/rvar"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BehSpawn struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity

	// SpawnInterval is the (possibly random) gap between spawn cycles, compiled
	// from an expression such as "5s" or "exp(30s)" and re-sampled each cycle.
	SpawnInterval       *rvar.Rvar
	SpawnOnInit         bool
	SpawnLimit          int
	SpawnMultiplication int
	SpawnActorUid       primitive.ObjectID

	ActiveActor     *entities.SimActor
	SpawnedCount    int
	RemainingSpawns int
	NextSpawnTime   tools.Seconds
}

func GetBehSpawn(entity *entities.SimEntity) *BehSpawn {
	return entity.Impl.(*BehSpawn)
}

func (self *BehSpawn) Step() {
	actualTime := self.World.Stepper.Now

	// If we still have an active actor that couldn't be passed yet, try to pass it first.
	if self.ActiveActor != nil {
		passed := self._TryToPassActor()
		if !passed {
			return
		}
	}

	// If we have leftover spawns from a previous multiplication cycle, continue them
	// before scheduling/processing the next interval-based spawn.
	if self.RemainingSpawns > 0 {
		self._ProcessRemainingSpawns()
		return
	}

	// No active actor, no remaining spawns from the previous cycle:
	// time to start a new spawn cycle when the interval elapses.
	if actualTime >= self.NextSpawnTime {
		self.Spawn()
	}
}

/*
Spawn tries to spawn entity, it has a limit of max spawn entities
*/
func (self *BehSpawn) Spawn() {
	if self.CanSpawn() {
		self.ForceSpawn()
	}
}

/*
ForceSpawn starts a new spawn cycle, attempting to spawn SpawnMultiplication entities.
Any that cannot be spawned/passed immediately are deferred via RemainingSpawns and
retried on subsequent steps.
*/
func (self *BehSpawn) ForceSpawn() {
	self.RemainingSpawns = self.SpawnMultiplication
	self._ProcessRemainingSpawns()
}

/*
_ProcessRemainingSpawns attempts to spawn and pass entities while there are remaining
slots in the current multiplication cycle. Stops as soon as one actor cannot be passed,
leaving it as ActiveActor so it will be retried on the next Step().
*/
func (self *BehSpawn) _ProcessRemainingSpawns() {
	for self.RemainingSpawns > 0 {

		// Respect the spawn limit during multiplication as well.
		if !self.CanSpawn() {
			self.RemainingSpawns = 0
			self._CalculateNextSpawn()
			return
		}

		self.ActiveActor = self.World.SpawnActorWithUid(self.SpawnActorUid)
		lgr.Ok("Spawned actor %v", self.ActiveActor)

		if self.ActiveActor == nil {
			// Spawning failed at the world level; abandon the cycle and reschedule.
			self.RemainingSpawns = 0
			self._CalculateNextSpawn()
			return
		}

		self.SpawnedCount++
		self.RemainingSpawns--

		self.World.UpdateActorState(self.Entity.RuntimeId, "count", self.SpawnedCount)

		self.ActiveActor.Position = self.Entity.Position
		self.NextSpawnTime = tools.MaxSeconds()

		passed := self._TryToPassActor()
		if !passed {
			// Couldn't deliver this actor to a connection — pause the cycle.
			// We'll retry on the next Step(); the remaining count is preserved.
			return
		}
	}

	// All multiplications for this cycle were spawned and passed successfully.
	self._CalculateNextSpawn()
}

/*
_TryToPassActor attempts to hand the current ActiveActor over to one of the entity's
connections. Returns true if the actor was successfully taken.
*/
func (self *BehSpawn) _TryToPassActor() bool {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {
			taken := connection.TakeActor(self.ActiveActor)
			if taken {
				self.ActiveActor = nil
				return true
			}
		}
	}

	return false
}

func (self *BehSpawn) _CalculateNextSpawn() {
	self.NextSpawnTime = self.World.Stepper.Now + tools.Seconds(self.SpawnInterval.GetInt64())

	self.World.CreateUpcomingEvent(simInternal.SimUpcomingEvent{
		Id:     self.Entity.RuntimeId,
		Type:   simInternal.UPC_EVNT_SPAWN,
		Second: self.NextSpawnTime,
	})

	lgr.Info(self.Entity.FormatInfo("Next spawn event at: (%d/%d)", self.NextSpawnTime, self.World.Stepper.End))
}

func (self *BehSpawn) CanSpawn() bool {
	if self.SpawnLimit <= 0 {
		return true
	} else {
		return self.SpawnedCount < self.SpawnLimit
	}
}
