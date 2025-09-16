package simulation

type SwitchStrategy int8

const (
	FIRST_FREE  SwitchStrategy = 0
	ROUND_ROBIN SwitchStrategy = 1
)

type SwitchBehaviour struct {
	World  *SimWorld
	Entity *SimEntity

	Mode SwitchStrategy

	RRIndex int
}

func NewSwitchBehaviour() ISimBehaviour {
	return &SwitchBehaviour{}
}

func (self *SwitchBehaviour) __Interface() {
	var _ ISimBehaviour = &SwitchBehaviour{}
	var _ ActorTakerBehaviour = &SwitchBehaviour{}
}

// Implementacia IBehaviour
func (self *SwitchBehaviour) SetWorld(world *SimWorld) {
	self.World = world
}

func (self *SwitchBehaviour) Init1() {
}

func (self *SwitchBehaviour) Init2() {
}

func (self *SwitchBehaviour) GetEntity() *SimEntity {
	return self.Entity
}

func (self *SwitchBehaviour) SetEntity(entity *SimEntity) {
	self.Entity = entity

	self.Mode = SwitchStrategy(entity.TypeData.GetInt8("mode", 0))
}

func (self *SwitchBehaviour) Step() {

}

// Take actor behaviour
func (self *SwitchBehaviour) TakeActor(actor *SimActor) bool {

	connections, hasConnections := self.World.SimConnections[self.Entity.Uid]

	if hasConnections {

		if self.Mode == FIRST_FREE {
			if hasConnections {
				for _, connection := range connections {

					taker, isITaker := connection.(ActorTakerBehaviour)

					if isITaker {
						taken := taker.TakeActor(actor)

						if taken {
							return true
						}
					}
				}
			}
		} else if self.Mode == ROUND_ROBIN {
			//pass

			conn := connections[self.RRIndex].(ActorTakerBehaviour)

			taken := conn.TakeActor(actor)

			if taken {
				self.RRIndex += 1

				if self.RRIndex >= len(connections) {
					self.RRIndex = 0
				}
				return true
			} else {
				return false
			}

		}
	}

	return false
}

func (self *SwitchBehaviour) CanTakeActor(actor *SimActor) bool {
	//TODO implement me
	panic("implement me")
}
