package simulation2

type BehavioursFactory struct {
	Behaviours map[string]func(entity *SimEntity)
}

var BEH_FACTORY = NewBehavioursFactory()

func NewBehavioursFactory() *BehavioursFactory {

	tmp := &BehavioursFactory{
		Behaviours: make(map[string]func(entity *SimEntity)),
	}

	tmp.Behaviours["buffer"] = NewBufferBehaviour
	tmp.Behaviours["process"] = NewProcessBehaviour
	tmp.Behaviours["spawn"] = NewSpawnBehaviour

	return tmp
}
