package simulation

func __CompilaitonCheck() {
	var _ ISimBehaviour = &ProcessBehaviour{}
	var _ ISimBehaviour = &SpawnBehaviour{}

	var _ ActorTakerBehaviour = &ProcessBehaviour{}

	var _ ActorTakerBehaviour = &BufferBehaviour{}
	var _ ActorProviderBehaviour = &BufferBehaviour{}

}
