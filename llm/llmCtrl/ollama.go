package llmCtrl

import (
	"fmt"
	"github.com/erikpa1/turtle/lg"
	"github.com/erikpa1/turtle/llm/llmModels"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type RunningModel struct {
	Port  int
	Model *llmModels.LLM
}

var OllamaModels = make(map[primitive.ObjectID]*RunningModel)
var OllamaModelsLock = sync.Mutex{}

func InitOllama() {

	localModels := ListLocalModels(primitive.ObjectID{})

	modelPort := 1143

	if len(localModels) > 0 {

		lg.LogI("Going to init Ollama")

		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {

			cmdString := fmt.Sprintf("set OLLAMA_HOST=localhost:%d && ollama serve &", modelPort)

			cmd = exec.Command(
				"cmd",
				"/C",
				cmdString,
			)

		} else {
			// Unix/Linux/Mac
			cmd = exec.Command(
				"sh",
				"-c",
				fmt.Sprintf("OLLAMA_HOST=localhost:%d ollama serve &", modelPort),
			)
		}

		output, err := cmd.Output()
		if err != nil {
			lg.LogE(err.Error())
		} else {
			lg.LogOk("LLama probably initialized [", output, "]")
		}
	}

	return

	if len(localModels) > 0 {

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

						cmdString := fmt.Sprintf("set OLLAMA_HOST=localhost:%d && ollama serve && ollama run %s", modelPort, model.ModelVersion)

						cmd = exec.Command(
							"cmd",
							"/C",
							cmdString,
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

func InstallOllama(cluster primitive.ObjectID, model string) string {
	if cluster.IsZero() {

		finalCommand := fmt.Sprintf("ollama pull %s", model)

		lg.LogOk("Final command: ", finalCommand)

		cmd := exec.Command(
			"sh",
			"-c",
			finalCommand,
		)
		// Capture the output
		output, err := cmd.Output()
		if err != nil {
			lg.LogE(err.Error())
		}
		return string(output)
	}
	//TODO naimplementovat RPC
	return ""
}

func OllamaStart() {
	cmd := exec.Command(
		"sh",
		"-c",
		"ollama serve &",
	)

	// Capture the output
	output, err := cmd.Output()

	if err != nil {
		lg.LogE(err.Error())
	} else {
		lg.LogOk(output)
	}

}

func OllamaList() string {
	cmd := exec.Command(
		"sh",
		"-c",
		"ollama list",
	)

	// Capture the output
	output, err := cmd.Output()

	if err != nil {
		lg.LogE(err.Error())
	}

	return modelsToHTML(parseOllamaOutput(string(output)))

}

type OllamaModel struct {
	Name     string
	ID       string
	Size     string
	Modified string
}

func parseOllamaOutput(output string) []OllamaModel {
	var models []OllamaModel

	lines := strings.Split(output, "\n")

	// Skip header line and empty lines
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "NAME") {
			continue
		}

		// Use regex to parse the line with proper spacing
		// This handles cases where fields might contain spaces
		re := regexp.MustCompile(`^(\S+(?:\s+\S+)*?)\s+([a-f0-9]{12})\s+([0-9.]+\s+[KMGT]B)\s+(.+)$`)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 5 {
			model := OllamaModel{
				Name:     strings.TrimSpace(matches[1]),
				ID:       strings.TrimSpace(matches[2]),
				Size:     strings.TrimSpace(matches[3]),
				Modified: strings.TrimSpace(matches[4]),
			}
			models = append(models, model)
		}
	}

	return models
}

func modelsToHTML(models []OllamaModel) string {
	if len(models) == 0 {
		return "<p>No models found.</p>"
	}

	html := `<table border="1" cellpadding="8" cellspacing="0" style="border-collapse: collapse; font-family: Arial, sans-serif;">
	<thead>
		<tr style="background-color: #f0f0f0;">
			<th>Name</th>
			<th>ID</th>
			<th>Size</th>
			<th>Modified</th>
		</tr>
	</thead>
	<tbody>`

	for i, model := range models {
		bgColor := ""
		if i%2 == 1 {
			bgColor = ` style="background-color: #f9f9f9;"`
		}

		html += fmt.Sprintf(`
		<tr%s>
			<td>%s</td>
			<td><code>%s</code></td>
			<td>%s</td>
			<td>%s</td>
		</tr>`, bgColor, model.Name, model.ID, model.Size, model.Modified)
	}

	html += `
	</tbody>
</table>`

	return html
}
