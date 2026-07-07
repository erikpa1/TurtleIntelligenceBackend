package behDelay

import (
	"slices"
	"turtle/simulation2/entities"

	"turtle/simulation/stats"
	"turtle/simulation2/rvar"
	"turtle/tools"
)

type BehDelay struct {
	World  *entities.SimWorld
	Entity *entities.SimEntity

	// Actors waiting to be released, keyed by their scheduled release time.
	// Multiple actors can share the same release time (they arrived in the
	// same simulation step), so the value is a slice preserving arrival order.
	Actors map[tools.Seconds][]*entities.SimActor

	// DelayTime is the (possibly random) delay each actor waits before being
	// released downstream. It is compiled from an expression such as "00:10",
	// "10s" or "exp(10s)" and sampled once per arriving actor.
	DelayTime *rvar.Rvar

	Statistics *stats.ProcessStats
}

func GetBehDelay(entity *entities.SimEntity) *BehDelay {
	return entities.CastImplementation[BehDelay](entity.Impl)
}

// TakeActor accepts an actor into the infinite delay buffer. It always returns
// true because the delay has no capacity limit. Each actor draws its own delay
// from DelayTime, so the wait can be fixed or random per actor.
func (self *BehDelay) TakeActor(actor *entities.SimActor) bool {
	now := self.World.Stepper.Now
	delay := tools.Seconds(self.DelayTime.GetInt64())
	releaseTime := now + delay

	self.Actors[releaseTime] = append(self.Actors[releaseTime], actor)
	actor.UpdatePosition(self.Entity.Position.RandomizeXZ(1))
	self.UpdateCountToClient()

	return true
}

// Step releases all actors whose delay has elapsed, passing them downstream one
// at a time. Any actor a downstream entity refuses is kept and retried next
// step. Because the world steps every entity on every tick, no explicit
// rescheduling is needed.
func (self *BehDelay) Step() {
	now := self.World.Stepper.Now

	// Collect release times that are due, sorted so we release in chronological
	// order (Go map iteration order is undefined).
	dueTimes := make([]tools.Seconds, 0, len(self.Actors))
	for t := range self.Actors {
		if t <= now {
			dueTimes = append(dueTimes, t)
		}
	}
	slices.Sort(dueTimes)

	connections := self.World.GetConnectionsOf(self.Entity.Uid)
	blocked := false

	for _, t := range dueTimes {
		bucket := self.Actors[t]
		remaining := bucket[:0]

		for _, actor := range bucket {
			if self._TryPassDownstream(actor, connections) {
				continue
			}
			remaining = append(remaining, actor)
		}

		if len(remaining) == 0 {
			delete(self.Actors, t)
		} else {
			// A blocked downstream means later actors must wait too; stop
			// pushing further buckets so order is preserved across ticks.
			self.Actors[t] = remaining
			blocked = true
			break
		}
	}

	if blocked {
		self.Statistics.BlockedTime += 1
	}

	self.UpdateCountToClient()
}

// _TryPassDownstream offers an actor to each connected entity in turn and
// returns true as soon as one accepts it.
func (self *BehDelay) _TryPassDownstream(actor *entities.SimActor, connections []*entities.SimEntity) bool {
	for _, conn := range connections {
		if conn.TakeActor(actor) {
			return true
		}
	}
	return false
}

// UpdateCountToClient publishes how many actors are currently waiting.
func (self *BehDelay) UpdateCountToClient() {
	self.World.UpdateActorState(self.Entity.RuntimeId, "count", self.WaitingCount())
}

// WaitingCount returns the total number of actors still inside the delay.
func (self *BehDelay) WaitingCount() int {
	total := 0
	for _, bucket := range self.Actors {
		total += len(bucket)
	}
	return total
}
