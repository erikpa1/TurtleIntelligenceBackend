package llmCtrl

//Source https://www.youtube.com/watch?v=1MZ2xb178NA
import (
	"turtle/lg"

	"github.com/tmc/langchaingo/prompts"
)

func RunStringPromptTemplate() {
	lg.LogI("Start lang chain templates")

	simpleTemplate := prompts.NewPromptTemplate(
		"Write a {{.content_type}} about {{.subject}}",
		[]string{"content_type", "subject"},
	)

	templateInput := map[string]any{
		"content_type": "poem",
		"subject":      "cats",
	}

	simple_prompt, _ := simpleTemplate.Format(templateInput)
	lg.LogE(simple_prompt)

}
