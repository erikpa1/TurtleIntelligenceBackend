package simulation2

type BehavioursFactory struct {
	Behaviours map[string]func(entity *SimEntity)
}

var BEH_FACTORY = NewBehavioursFactory()

func NewBehavioursFactory() *BehavioursFactory {

	tmp := &BehavioursFactory{
		Behaviours: make(map[string]func(entity *SimEntity)),
	}

	tmp.Behaviours["process"] = NewProcessBehaviour
	tmp.Behaviours["spawn"] = NewSpawnBehaviour
	tmp.Behaviours["buffer"] = NewBufferBehaviour
	tmp.Behaviours["sink"] = NewSinkBehaviour

	return tmp
}
