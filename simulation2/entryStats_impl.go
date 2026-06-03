package simulation2

import "turtle/tools"

type EntryRecord struct {
	ActorId int64         // adjust type to match SimActor.Uid
	Time    tools.Seconds // adjust type to match Stepper.Now
}

type BehEntryStatistics struct {
	World  *SimWorld
	Entity *SimEntity

	// Statistics
	Entries    []EntryRecord // chronological log of every actor that entered
	EntryCount int           // total number of actors that entered
	LastEntry  tools.Seconds // time of the most recent entry (adjust type)
	FirstEntry tools.Seconds // time of the first entry (adjust type)
}

func GetBehEntryStatistics(entity *SimEntity) *BehEntryStatistics {
	return entity.Impl.(*BehEntryStatistics)
}

func (self *BehEntryStatistics) TakeActor(actor *SimActor) bool {

	now := self.World.Stepper.Now

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {
		for _, connection := range connections {
			taken := connection.TakeActor(actor)

			if taken {
				self.recordEntry(actor, now)
				return true
			}
		}
	}

	return false
}

func (self *BehEntryStatistics) recordEntry(actor *SimActor, now tools.Seconds) {
	if self.EntryCount == 0 {
		self.FirstEntry = now
	}

	self.Entries = append(self.Entries, EntryRecord{
		ActorId: actor.Id,
		Time:    now,
	})

	self.EntryCount++
	self.LastEntry = now
}
