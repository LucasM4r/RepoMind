package rag

import (
	"regexp"
	"sort"
	"strings"

	"github.com/LucasM4r/repomind/internal/domain"
)

var tokenizeRegex = regexp.MustCompile("[^a-z0-9]+")

func tokenize(text string) []string {
	text = strings.ToLower(text)
	text = tokenizeRegex.ReplaceAllString(text, " ")
	return strings.Fields(text)
}

func computeLexicalScore(prompt, chunkContent string) float32 {
	promptTokens := tokenize(prompt)
	chunkTokens := tokenize(chunkContent)
	if len(promptTokens) == 0 || len(chunkTokens) == 0 {
		return 0.0
	}

	// Calculate the Jaccard similarity index to measure
	// the lexical similarity between the prompt and chunk content
	// Jaccard index = |Intersection(A, B)| / |Union(A, B)|
	chunkSet := make(map[string]struct{}, len(chunkTokens))
	for _, token := range chunkTokens {
		chunkSet[token] = struct{}{}
	}

	promptSet := make(map[string]struct{}, len(promptTokens))
	for _, token := range promptTokens {
		promptSet[token] = struct{}{}
	}

	intersection := 0
	for token := range promptSet {
		if _, exists := chunkSet[token]; exists {
			intersection++
		}
	}

	union := len(promptSet) + len(chunkSet) - intersection
	if union == 0 {
		return 0.0
	}

	return float32(intersection) / float32(union)
}

func rerankChunks(prompt string, chunks []domain.RetrievedChunk, finalTopK int) []domain.RetrievedChunk {
	if len(chunks) == 0 || finalTopK <= 0 {
		return []domain.RetrievedChunk{}
	}

	type scoredChunk struct {
		chunk        domain.RetrievedChunk
		vectorScore  float32
		lexicalScore float32
		finalScore   float32
	}

	scoredChunks := make([]scoredChunk, len(chunks))

	for i, chunk := range chunks {
		vectorScore := float32(1.0) - chunk.Distance
		if vectorScore < 0 {
			vectorScore = 0
		}
		if vectorScore > 1 {
			vectorScore = 1
		}

		lexicalScore := computeLexicalScore(prompt, chunk.Content)
		finalScore := 0.7*vectorScore + 0.3*lexicalScore

		scoredChunks[i] = scoredChunk{
			chunk:        chunk,
			vectorScore:  vectorScore,
			lexicalScore: lexicalScore,
			finalScore:   finalScore,
		}
	}

	sort.Slice(scoredChunks, func(i, j int) bool {
		if scoredChunks[i].finalScore == scoredChunks[j].finalScore {
			return scoredChunks[i].chunk.Distance < scoredChunks[j].chunk.Distance
		}
		return scoredChunks[i].finalScore > scoredChunks[j].finalScore
	})

	if finalTopK > len(scoredChunks) {
		finalTopK = len(scoredChunks)
	}

	result := make([]domain.RetrievedChunk, finalTopK)
	for i := 0; i < finalTopK; i++ {
		result[i] = scoredChunks[i].chunk
	}

	return result
}
