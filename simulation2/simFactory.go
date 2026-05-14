package simulation2

var BEH_FACTORY = NewBehavioursFactory()

type BehavioursFactory struct {
	Behaviours map[string]func(entity *SimEntity)
}

func NewBehavioursFactory() *BehavioursFactory {

	tmp := &BehavioursFactory{
		Behaviours: make(map[string]func(entity *SimEntity)),
	}

	tmp.Behaviours["buffer"] = NewBufferBehaviour
	tmp.Behaviours["process"] = NewProcessBehaviour
	tmp.Behaviours["process"] = NewSpawnBehaviour

	return tmp
}
