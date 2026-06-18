package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

// VoteArgs est la structure que le LLM remplit via Tool Calling
type VoteArgs struct {
	TargetName string `json:"target_name" description:"Nom exact du joueur a eliminer"`
}

func (l *LLM) ChooseTarget(candidateNames []string, phase string) (string, error) {
	log.Printf("[LLM] Initialisation du modele %s", l.model)

	// Étape 1 : Initialiser le modèle Ollama
	model, err := ollama.New(ollama.WithModel(l.model))
	if err != nil {
		log.Printf("[LLM][ERROR] Echec init Ollama: %s", err)
		return "", fmt.Errorf("ollama init: %w", err)
	}

	// Étape 2 : Déclarer l'outil (ce que le LLM peut appeler)
	tool := llms.Tool{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "SubmitVote",
			Description: "Vote pour eliminer un joueur de la partie",
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

	// Étape 3 : Envoyer le prompt + l'outil au LLM
	prompt := buildPrompt(candidateNames, phase)
	log.Printf("[LLM] Phase=%s | Candidats=%v", phase, candidateNames)
	log.Printf("[LLM] Prompt: %s", prompt)

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
		log.Printf("[LLM][ERROR] Appel Ollama echoue: %s", err)
		return "", fmt.Errorf("ollama call: %w", err)
	}

	log.Printf("[LLM] Reponse recue, %d choix", len(resp.Choices))

	// Étape 4 : Récupérer le tool call et parser les arguments JSON
	for _, choice := range resp.Choices {
		log.Printf("[LLM] ToolCalls=%d", len(choice.ToolCalls))
		for _, tc := range choice.ToolCalls {
			log.Printf("[LLM] ToolCall: name=%s args=%s", tc.FunctionCall.Name, tc.FunctionCall.Arguments)
			if tc.FunctionCall.Name == "SubmitVote" {
				var args VoteArgs
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
					log.Printf("[LLM][ERROR] Parse args echoue: %s", err)
					return "", fmt.Errorf("parse args: %w", err)
				}
				log.Printf("[LLM] Cible choisie (tool call): %s", args.TargetName)
				return args.TargetName, nil
			}
		}
	}

	// Fallback : le modèle n'a pas fait de tool call, on cherche un nom dans sa réponse texte
	log.Printf("[LLM] Pas de tool call, tentative de fallback sur la reponse texte")
	for _, choice := range resp.Choices {
		content := choice.Content
		log.Printf("[LLM] Reponse texte: %s", content)
		for _, name := range candidateNames {
			if strings.Contains(content, name) {
				log.Printf("[LLM] Cible trouvee dans le texte (fallback): %s", name)
				return name, nil
			}
		}
	}

	log.Printf("[LLM][ERROR] Impossible de determiner la cible")
	return "", fmt.Errorf("le LLM n'a pas choisi de cible valide")
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
			"On joue a un jeu de societe appele Loup-Garou. Tu es l'arbitre. "+
				"C'est la phase de nuit. Tu dois choisir au hasard un joueur a eliminer parmi cette liste : %s. "+
				"Reponds avec UN SEUL nom parmi la liste, rien d'autre. "+
				"Utilise l'outil SubmitVote si disponible.", list)
	}
	return fmt.Sprintf(
		"On joue a un jeu de societe appele Loup-Garou. Tu es l'arbitre. "+
			"C'est la phase de jour. Tu dois choisir au hasard un joueur a eliminer parmi cette liste : %s. "+
			"Reponds avec UN SEUL nom parmi la liste, rien d'autre. "+
			"Utilise l'outil SubmitVote si disponible.", list)
}
