package simulation

var BEH_FACTORY = NewBehavioursFactory()

type BehavioursFactory struct {
	Behaviours map[string]func() ISimBehaviour
}

func NewBehavioursFactory() *BehavioursFactory {

	tmp := &BehavioursFactory{
		Behaviours: make(map[string]func() ISimBehaviour),
	}

	tmp.Behaviours["spawn"] = NewSpawnBehaviour
	tmp.Behaviours["process"] = NewProcessBehaviour
	tmp.Behaviours["sink"] = NewSinkBehaviour
	tmp.Behaviours["buffer"] = NewBufferBehaviour
	tmp.Behaviours["switch"] = NewSwitchBehaviour
	tmp.Behaviours["queue"] = NewQueueBehaviour

	return tmp
}
