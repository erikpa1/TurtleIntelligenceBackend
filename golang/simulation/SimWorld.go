package simulation

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/lg"
	"turtle/modelsApp"
	"turtle/server"
	"turtle/tools"
)

type SimWorld struct {
	Uid      primitive.ObjectID
	Name     string
	Entities map[primitive.ObjectID]*modelsApp.Entity
	Stepper  SimStepper
	IsOnline bool
}

func NewSimWorld() *SimWorld {
	tmp := &SimWorld{}
	tmp.Stepper.End = 100
	tmp.IsOnline = true

	return tmp

}

func (self *SimWorld) RunSimulation() {

	var second tools.Seconds = 0

	for second = 0; second < self.Stepper.End; second++ {
		self.Step()
		self.Stepper.Step()

		server.MYIO.EmitSync("simstep", bson.M{
			"second": second,
		})
	}

}

func (self *SimWorld) Step() {

	lg.LogI(fmt.Sprintf("Step (%d/%d)", self.Stepper.Now, self.Stepper.End))

}
