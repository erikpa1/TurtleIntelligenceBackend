package simulation2

import (
	"turtle/simulation/stats"
	"turtle/tools"
)

type BehDelay struct {
	World  *SimWorld
	Entity *SimEntity

	// Actors waiting to be released, keyed by their scheduled release time.
	// Multiple actors can share the same release time (they arrived in the
	// same simulation step), so the value is a slice preserving arrival order.
	Actors map[tools.Seconds][]*SimActor

	// Configured delay duration (e.g. "00:10") parsed into seconds on use.
	DelayTime string

	Statistics *stats.ProcessStats
}

func GetBehDelay(entity *SimEntity) *BehDelay {
	return CastImplementation[BehDelay](entity.Impl)
}

// TakeActor accepts an actor into the infinite delay buffer.
// It always returns true because the delay has no capacity limit.
func (d *BehDelay) TakeActor(actor *SimActor) bool {
	delay := tools.ParseSeconds(d.DelayTime)
	releaseTime := d.World.Now() + delay

	d.Actors[releaseTime] = append(d.Actors[releaseTime], actor)

	d.Statistics.OnEnter(d.World.Now())

	// Make sure we get stepped at the moment this actor is due to leave.
	d.Entity.ScheduleStep(releaseTime)

	return true
}

// Step releases all actors whose delay has elapsed. Actors are passed
// downstream one at a time; any actor that the downstream refuses to
// take is kept in the buffer and will be retried on the next step.
func (d *BehDelay) Step() {
	now := d.World.Now()

	// Collect the release times that are due, so we can iterate in
	// chronological order (Go maps have no defined iteration order).
	dueTimes := make([]tools.Seconds, 0, len(d.Actors))
	for t := range d.Actors {
		if t <= now {
			dueTimes = append(dueTimes, t)
		}
	}
	tools.SortSeconds(dueTimes)

	for _, t := range dueTimes {
		bucket := d.Actors[t]
		remaining := bucket[:0]

		for _, actor := range bucket {
			if d.Entity.OfferActorDownstream(actor) {
				d.Statistics.OnLeave(now)
			} else {
				// Downstream refused — keep this actor (and any after
				// it, to preserve FIFO order within the bucket).
				remaining = append(remaining, actor)
			}
		}

		if len(remaining) == 0 {
			delete(d.Actors, t)
		} else {
			d.Actors[t] = remaining
			// Stop pushing further buckets: a blocked downstream means
			// later actors must wait too, preserving order.
			break
		}
	}

	// If anything is still waiting, make sure we get stepped again.
	if len(d.Actors) > 0 {
		d.Entity.ScheduleStep(now + 1)
	}
}
