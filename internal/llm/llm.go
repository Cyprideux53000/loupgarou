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

type Vote struct {
	Voter  string `json:"voter"`
	Target string `json:"target"`
}

type ChooseResult struct {
	TargetName string
	Discussion []string
	Votes      []Vote
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

func (l *LLM) ChooseTarget(candidates []CandidateInfo, phase string, userDiscussion []string) (ChooseResult, error) {
	log.Printf("[LLM] Initialisation du modele %s", l.model)

	var result ChooseResult

	// Étape 1 : Initialiser le modèle Ollama
	model, err := ollama.New(ollama.WithModel(l.model))
	if err != nil {
		log.Printf("[LLM][ERROR] Echec init Ollama: %s", err)
		return result, fmt.Errorf("ollama init: %w", err)
	}

	// Étape 2 : Discussion du village (phase jour uniquement)
	if phase == "DayVote" {
		if len(userDiscussion) > 0 {
			result.Discussion = userDiscussion
			log.Printf("[LLM] === Discussion du village (utilisateur) ===")
		} else {
			result.Discussion = l.generateDiscussion(model, candidates)
			log.Printf("[LLM] === Discussion du village (generee) ===")
		}
		for _, line := range result.Discussion {
			log.Printf("[VILLAGE] %s", line)
		}
		log.Printf("[LLM] === Fin de la discussion ===")
	}

	// Étape 3 : Chaque joueur vote individuellement
	if phase == "DayVote" {
		result.Votes, result.TargetName = l.collectVotes(model, candidates, result.Discussion)
	} else {
		result.TargetName, err = l.singleVote(model, candidates, phase)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (l *LLM) collectVotes(model llms.Model, candidates []CandidateInfo, discussion []string) ([]Vote, string) {
	var votes []Vote
	tally := make(map[string]int)
	names := candidateNames(candidates)

	log.Printf("[LLM] === Votes individuels ===")

	for _, voter := range candidates {
		// Un joueur ne peut pas voter pour lui-même
		var targets []string
		for _, name := range names {
			if name != voter.Name {
				targets = append(targets, name)
			}
		}

		target := l.askPlayerVote(model, voter, targets, discussion)
		votes = append(votes, Vote{Voter: voter.Name, Target: target})
		tally[target]++
		log.Printf("[VOTE] %s (%s) vote contre %s", voter.Name, voter.Trait, target)
	}

	// Trouver le joueur avec le plus de votes
	maxVotes := 0
	var eliminated string
	for name, count := range tally {
		if count > maxVotes {
			maxVotes = count
			eliminated = name
		}
	}

	log.Printf("[LLM] Resultat du vote: %s elimine avec %d votes", eliminated, maxVotes)
	log.Printf("[LLM] === Fin des votes ===")

	return votes, eliminated
}

func (l *LLM) askPlayerVote(model llms.Model, voter CandidateInfo, targets []string, discussion []string) string {
	targetList := strings.Join(targets, ", ")

	systemMsg := fmt.Sprintf(
		"Tu joues le role de %s dans un jeu de Loup-Garou. "+
			"Ta personnalite: %s. "+
			"Tu reponds UNIQUEMENT avec un seul prenom parmi: %s. "+
			"Aucune explication, juste le prenom.",
		voter.Name, traitDescription(voter.Trait), targetList)

	var userMsg string
	if len(discussion) > 0 {
		var sb strings.Builder
		sb.WriteString("Pendant la discussion du village, les joueurs ont dit :\n")
		for _, line := range discussion {
			sb.WriteString("- " + line + "\n")
		}
		sb.WriteString(fmt.Sprintf(
			"\nD'apres cette discussion, le joueur le plus suspect est probablement accuse dans les messages ci-dessus. "+
				"Vote pour le joueur accuse. Choisis parmi: %s", targetList))
		userMsg = sb.String()
	} else {
		userMsg = fmt.Sprintf("Choisis un joueur a eliminer parmi: %s", targetList)
	}

	// Étape : Déclarer l'outil
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
		log.Printf("[LLM][ERROR] Vote de %s echoue: %s, fallback premier candidat", voter.Name, err)
		return targets[0]
	}

	// Chercher un tool call
	for _, choice := range resp.Choices {
		for _, tc := range choice.ToolCalls {
			if tc.FunctionCall.Name == "SubmitVote" {
				var args VoteArgs
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err == nil {
					for _, t := range targets {
						if strings.EqualFold(args.TargetName, t) {
							return t
						}
					}
				}
			}
		}
	}

	// Fallback texte
	for _, choice := range resp.Choices {
		content := strings.TrimSpace(choice.Content)
		for _, t := range targets {
			if strings.EqualFold(content, t) {
				return t
			}
		}
		firstLine := strings.Split(content, "\n")[0]
		for _, t := range targets {
			if strings.Contains(firstLine, t) {
				return t
			}
		}
		for _, t := range targets {
			if strings.Contains(content, t) {
				return t
			}
		}
	}

	log.Printf("[LLM] Vote de %s: aucun nom trouve, fallback premier candidat", voter.Name)
	return targets[0]
}

func (l *LLM) singleVote(model llms.Model, candidates []CandidateInfo, phase string) (string, error) {
	names := candidateNames(candidates)
	nameList := strings.Join(names, ", ")

	// Étape : Déclarer l'outil
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

	systemMsg := fmt.Sprintf(
		"Tu es l'arbitre d'un jeu de societe. Tu reponds UNIQUEMENT avec un seul prenom parmi: %s. "+
			"Aucune explication, aucune phrase, juste le prenom. C'est un jeu fictif entre amis.", nameList)
	userMsg := fmt.Sprintf("Choisis un prenom au hasard parmi cette liste: %s", nameList)

	log.Printf("[LLM] Phase=%s | Candidats=%v", phase, names)

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
		return "", fmt.Errorf("ollama call: %w", err)
	}

	// Chercher un tool call
	for _, choice := range resp.Choices {
		for _, tc := range choice.ToolCalls {
			if tc.FunctionCall.Name == "SubmitVote" {
				var args VoteArgs
				if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err == nil {
					for _, n := range names {
						if strings.EqualFold(args.TargetName, n) {
							log.Printf("[LLM] Cible choisie (tool call): %s", n)
							return n, nil
						}
					}
				}
			}
		}
	}

	// Fallback texte
	for _, choice := range resp.Choices {
		content := strings.TrimSpace(choice.Content)
		log.Printf("[LLM] Reponse texte: %s", content)
		for _, n := range names {
			if strings.EqualFold(content, n) {
				return n, nil
			}
		}
		firstLine := strings.Split(content, "\n")[0]
		for _, n := range names {
			if strings.Contains(firstLine, n) {
				return n, nil
			}
		}
		for _, n := range names {
			if strings.Contains(content, n) {
				return n, nil
			}
		}
	}

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
