package agents

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pipeline struct {
	Steps []*PipelineStep `json:"steps"`
}

type PipelineStep struct {
	Index     int                `json:"index"`
	Name      string             `json:"name"`
	NodeUid   primitive.ObjectID `json:"nodeUid"`
	Status    string             `json:"status"`
	StartedAt tools.Milliseconds `json:"startedAt"`
	EndedAt   tools.Milliseconds `json:"endedAt"`
	Duration  tools.Milliseconds `json:"duration"`
	HelpData  string             `json:"helpData"`
}

func (self *Pipeline) AddStep(step *PipelineStep) int {
	currentIndex := len(self.Steps)
	step.Index = currentIndex

	self.Steps = append(self.Steps, step)
	return currentIndex
}

func (self *Pipeline) NewStep() *PipelineStep {
	step := &PipelineStep{}
	step.Index = len(self.Steps)

	self.Steps = append(self.Steps, step)
	return step
}
func (self *Pipeline) NewStepFromNode(node *LLMAgentNode) *PipelineStep {
	step := self.NewStep()
	step.Name = node.Name
	step.NodeUid = node.Uid
	return step
}

func (self *PipelineStep) Start() {
	self.Status = "started"
	self.StartedAt = tools.GetTimeNowMillis()
}

func (self *PipelineStep) End() {
	self.Status = "end"
	self.EndedAt = tools.GetTimeNowMillis()
	self.Duration = self.EndedAt - self.StartedAt
}
