package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type LLM struct {
	model string
}

func New(model string) *LLM {
	return &LLM{model: model}
}

func (l *LLM) GetModel() string {
	return l.model
}

func (l *LLM) ChooseTarget(candidateNames []string, phase string) (string, error) {
	model, err := ollama.New(ollama.WithModel(l.model))
	if err != nil {
		return "", fmt.Errorf("ollama init: %w", err)
	}

	tool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "SubmitVote",
			Description: "Vote pour eliminer un joueur",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"target_name": map[string]any{
						"type":        "string",
						"description": "Nom exact du joueur a eliminer",
					},
				},
				"required": []string{"target_name"},
			},
		},
	}

	prompt := buildPrompt(candidateNames, phase)

	resp, err := model.GenerateContent(
		context.Background(),
		[]llms.MessageContent{
			{Role: llms.ChatMessageTypeHuman, Parts: []llms.ContentPart{
				llms.TextContent{Text: prompt},
			}},
		},
		llms.WithTools([]llms.Tool{tool}),
	)
	if err != nil {
		return "", fmt.Errorf("ollama call: %w", err)
	}

	for _, choice := range resp.Choices {
		for _, tc := range choice.ToolCalls {
			if tc.FunctionCall.Name == "SubmitVote" {
				var args struct {
					TargetName string `json:"target_name"`
				}
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
					return "", fmt.Errorf("parse args: %w", err)
				}
				return args.TargetName, nil
			}
		}
	}

	return "", fmt.Errorf("le LLM n'a pas utilise l'outil SubmitVote")
}

func buildPrompt(names []string, phase string) string {
	list := ""
	for i, n := range names {
		if i > 0 {
			list += ", "
		}
		list += n
	}

	if phase == "wolfAttack" {
		return fmt.Sprintf(
			"Tu es le meneur d'un jeu de Loup-Garou. C'est la nuit. "+
				"Les loups doivent eliminer un villageois parmi : %s. "+
				"Utilise l'outil SubmitVote pour voter.", list)
	}
	return fmt.Sprintf(
		"Tu es le meneur d'un jeu de Loup-Garou. C'est le jour. "+
			"Le village doit eliminer un suspect parmi : %s. "+
			"Utilise l'outil SubmitVote pour voter.", list)
}
