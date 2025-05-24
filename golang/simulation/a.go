package simulation

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
	"turtle/ctrlApp"
	"turtle/lg"
	"turtle/server"
	"turtle/tools"
)

type _RunningSim struct {
	Uid   primitive.ObjectID
	Model primitive.ObjectID
	User  primitive.ObjectID
}

var RUNNING_SIMS = make(map[primitive.ObjectID]*_RunningSim)
var RUNNING_SIMS_LOCK = new(sync.Mutex)

func RunSimulation(modelUid primitive.ObjectID, simParams bson.M) {

	entities := ctrlApp.QueryWorldEntities(bson.M{"model": modelUid})
	connections := ctrlApp.ListConnectionsOfWorld(modelUid)

	ctx := context.Background()

	for _, entity := range entities {
		lg.LogI(entity)
	}

	world := NewSimWorld()
	world.IsOnline = true
	world.LoadEntities(entities)
	world.LoadConnections(connections)
	world.PrepareSimulation()

	runSim := &_RunningSim{}
	runSim.Uid = primitive.NewObjectID()
	runSim.Model = modelUid

	RUNNING_SIMS_LOCK.Lock()
	RUNNING_SIMS[runSim.Uid] = &_RunningSim{}
	RUNNING_SIMS_LOCK.Unlock()

	go func() {
		tools.Recover("Failed to run simulation")

		var second tools.Seconds = 0

	simulationLoop:
		for second = 0; second < world.Stepper.End; second++ {

			select {
			case <-ctx.Done():
				{
					lg.LogW("User canceled simulation")
					break simulationLoop
				}
			default:
				{
					// Record start time before world.Step()
					stepStart := time.Now()

					world.Step()
					world.Stepper.Step()

					server.MYIO.EmitSync("simstep", bson.M{
						"second":    second,
						"spawned":   []bson.M{},
						"unspawned": []int64{},
						"states":    []bson.M{},
					})

					world.ClearStates()

					// Calculate elapsed time for the step
					stepDuration := time.Since(stepStart)

					// Calculate remaining sleep time (1 second target - step duration)
					targetDuration := 1000 * time.Millisecond
					remainingSleep := targetDuration - stepDuration

					// Only sleep if there's remaining time
					if remainingSleep > 0 {
						time.Sleep(remainingSleep)
					}
				}
			}
		}

		RUNNING_SIMS_LOCK.Lock()
		delete(RUNNING_SIMS, runSim.Uid)
		RUNNING_SIMS_LOCK.Unlock()
	}()
}
