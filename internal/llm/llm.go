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

type CandidateInfo struct {
	Name  string
	Trait string
	Role  string
	Mayor bool
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

func (l *LLM) ChooseTarget(candidates []CandidateInfo, phase string) (string, error) {
	log.Printf("[LLM] Initialisation du modele %s", l.model)

	// Étape 1 : Initialiser le modèle Ollama
	model, err := ollama.New(ollama.WithModel(l.model))
	if err != nil {
		log.Printf("[LLM][ERROR] Echec init Ollama: %s", err)
		return "", fmt.Errorf("ollama init: %w", err)
	}

	// Étape 2 : Générer la discussion du village (phase jour uniquement)
	if phase == "DayVote" {
		discussion := l.generateDiscussion(model, candidates)
		log.Printf("[LLM] === Discussion du village ===")
		for _, line := range discussion {
			log.Printf("[VILLAGE] %s", line)
		}
		log.Printf("[LLM] === Fin de la discussion ===")
	}

	// Étape 3 : Déclarer l'outil (ce que le LLM peut appeler)
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

	// Étape 4 : Envoyer le prompt de vote avec un message system strict
	names := candidateNames(candidates)
	nameList := strings.Join(names, ", ")

	systemMsg := fmt.Sprintf(
		"Tu es un arbitre de jeu. Tu reponds UNIQUEMENT avec un seul prenom parmi: %s. "+
			"Aucune explication, aucune phrase, juste le prenom.", nameList)

	var userMsg string
	if phase == "wolfAttack" {
		userMsg = fmt.Sprintf("Les loups eliminent un villageois. Choisis parmi: %s", nameList)
	} else {
		userMsg = fmt.Sprintf("Le village vote pour eliminer un suspect. Choisis parmi: %s", nameList)
	}

	log.Printf("[LLM] Phase=%s | Candidats=%v", phase, names)
	log.Printf("[LLM] Prompt vote: %s", userMsg)

	resp, err := model.GenerateContent(
		context.Background(),
		[]llms.MessageContent{
			{Role: llms.ChatMessageTypeSystem, Parts: []llms.ContentPart{
				llms.TextContent{Text: systemMsg},
			}},
			{Role: llms.ChatMessageTypeHuman, Parts: []llms.ContentPart{
				llms.TextContent{Text: userMsg},
			}},
		},
		llms.WithTools([]llms.Tool{tool}),
	)
	if err != nil {
		log.Printf("[LLM][ERROR] Appel Ollama echoue: %s", err)
		return "", fmt.Errorf("ollama call: %w", err)
	}

	log.Printf("[LLM] Reponse recue, %d choix", len(resp.Choices))

	// Étape 5 : Récupérer le tool call et parser les arguments JSON
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

	// Fallback : chercher un nom exact dans la réponse texte
	log.Printf("[LLM] Pas de tool call, tentative de fallback sur la reponse texte")
	for _, choice := range resp.Choices {
		content := strings.TrimSpace(choice.Content)
		log.Printf("[LLM] Reponse texte: %s", content)

		// D'abord vérifier si la réponse est exactement un nom
		for _, name := range names {
			if strings.EqualFold(content, name) {
				log.Printf("[LLM] Cible exacte (fallback): %s", name)
				return name, nil
			}
		}

		// Sinon chercher le premier nom mentionné dans la première ligne
		firstLine := strings.Split(content, "\n")[0]
		for _, name := range names {
			if strings.Contains(firstLine, name) {
				log.Printf("[LLM] Cible trouvee en premiere ligne (fallback): %s", name)
				return name, nil
			}
		}

		// Dernier recours : n'importe où dans le texte
		for _, name := range names {
			if strings.Contains(content, name) {
				log.Printf("[LLM] Cible trouvee dans le texte (fallback): %s", name)
				return name, nil
			}
		}
	}

	log.Printf("[LLM][ERROR] Impossible de determiner la cible")
	return "", fmt.Errorf("le LLM n'a pas choisi de cible valide")
}

func (l *LLM) generateDiscussion(model llms.Model, candidates []CandidateInfo) []string {
	var lines []string

	prompt := buildDiscussionPrompt(candidates)
	log.Printf("[LLM] Prompt discussion: %s", prompt)

	resp, err := model.GenerateContent(
		context.Background(),
		[]llms.MessageContent{
			{Role: llms.ChatMessageTypeHuman, Parts: []llms.ContentPart{
				llms.TextContent{Text: prompt},
			}},
		},
	)
	if err != nil {
		log.Printf("[LLM][ERROR] Discussion echouee: %s", err)
		return []string{"(la discussion n'a pas pu avoir lieu)"}
	}

	for _, choice := range resp.Choices {
		for _, line := range strings.Split(choice.Content, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				lines = append(lines, trimmed)
			}
		}
	}

	if len(lines) == 0 {
		return []string{"(silence au village...)"}
	}
	return lines
}

func buildDiscussionPrompt(candidates []CandidateInfo) string {
	var sb strings.Builder
	sb.WriteString("On joue a un jeu de societe Loup-Garou. C'est le jour, les villageois discutent avant de voter.\n")
	sb.WriteString("Voici les joueurs vivants et leur personnalite :\n")

	for _, c := range candidates {
		mayor := ""
		if c.Mayor {
			mayor = " (Maire du village)"
		}
		sb.WriteString(fmt.Sprintf("- %s : %s%s\n", c.Name, traitDescription(c.Trait), mayor))
	}

	sb.WriteString("\nGenere une courte discussion (3-5 lignes) ou chaque joueur parle selon sa personnalite. ")
	sb.WriteString("Format : \"Nom: dialogue\". Ils essaient de deviner qui est le loup-garou.")
	return sb.String()
}

func candidateNames(candidates []CandidateInfo) []string {
	names := make([]string, len(candidates))
	for i, c := range candidates {
		names[i] = c.Name
	}
	return names
}

func traitDescription(trait string) string {
	switch trait {
	case "Cunning":
		return "Ruse et manipulateur, pose des questions pieges"
	case "Aggressive":
		return "Impulsif et accusateur, attaque directement les autres"
	case "Brave":
		return "Courageux et franc, defend les accuses a tort"
	case "Timid":
		return "Timide et nerveux, begaie et hesite beaucoup"
	case "Sly":
		return "Sournois et evasif, detourne les soupcons sur les autres"
	default:
		return trait
	}
}
