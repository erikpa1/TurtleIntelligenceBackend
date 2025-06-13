package llmCtrl

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"turtle/credentials"
	"turtle/lg"
	"turtle/llm/llmModels"
)

type RunningModel struct {
	Port  int
	Model *llmModels.LlmModel
}

var OllamaModels = make(map[primitive.ObjectID]*RunningModel)
var OllamaModelsLock = sync.Mutex{}

func InitOllama() {

	if credentials.IsSlaveApplication() {

		organization := primitive.ObjectID{}

		lg.LogOk("Slave, going to init agents")

		modelPort := 11434

		//TODO toto nie su LLM clustre, ale clustre ako take

		for _, cluster := range ListLLMClusters(organization) {

			if strings.Contains(cluster.Url, "localhost") {

				for _, model := range ListLLMModels(cluster.Org) {

					lg.LogI("Going to load model: ", model.ModelVersion, "on port: ", modelPort)

					var cmd *exec.Cmd

					if runtime.GOOS == "windows" {

						cmd = exec.Command(
							"cmd",
							"/C",
							fmt.Sprintf("set OLLAMA_HOST=localhost:%d && ollama serve && ollama run %s", modelPort, model.ModelVersion),
						)

					} else {
						// Unix/Linux/Mac
						cmd = exec.Command(
							"sh",
							"-c",
							fmt.Sprintf("OLLAMA_HOST=localhost:%d ollama serve && ollama run %s", modelPort, model.ModelVersion),
						)
					}

					err := cmd.Run()
					if err != nil {
						lg.LogE(err.Error())
					} else {
						OllamaModelsLock.Lock()
						tmp := RunningModel{}
						tmp.Model = model
						tmp.Port = modelPort
						OllamaModels[organization] = &tmp
						OllamaModelsLock.Unlock()
					}
				}
			}

		}

	}

}

func ListAndCheckRunningOllamas() []string {
	return []string{}
}
