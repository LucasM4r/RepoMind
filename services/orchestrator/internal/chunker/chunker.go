package chunker

import "strings"

type Chunk struct {
	Owner    string
	Repo     string
	Filename string
	Content  string
	Size     int
}

type Chunker interface {
	Split(owner, repo, filename, text string) []Chunk
}

type CharChunker struct {
	MaxChunkSize int
	Overlap      int
}

type LineChunker struct {
	MaxChunkSize int
	OverlapLines int
}

func NewCharChunker(maxChunkSize, overlap int) *CharChunker {
	return &CharChunker{
		MaxChunkSize: maxChunkSize,
		Overlap:      overlap,
	}
}

func (c *CharChunker) Split(owner, repo, filename, text string) []Chunk {
	runes := []rune(text)
	var chunks []Chunk
	length := len(runes)

	step := c.MaxChunkSize - c.Overlap
	if step <= 0 {
		step = 1
	}

	for i := 0; i < length; i += step {
		end := i + c.MaxChunkSize
		if end > length {
			end = length
		}

		content := string(runes[i:end])
		chunks = append(chunks, Chunk{
			Owner:    owner,
			Repo:     repo,
			Filename: filename,
			Content:  content,
			Size:     len(content),
		})
	}

	return chunks
}

func NewLineChunker(maxChunkSize, overlapLines int) *LineChunker {
	return &LineChunker{
		MaxChunkSize: maxChunkSize,
		OverlapLines: overlapLines,
	}
}

func (c *LineChunker) Split(owner, repo, filename, text string) []Chunk {
	lines := strings.Split(text, "\n")
	var chunks []Chunk
	var currentChunk []string
	var overlapBuffer []string
	currentLength := 0

	for _, line := range lines {
		lineLen := len(line) + 1

		if currentLength+lineLen > c.MaxChunkSize && len(currentChunk) > 0 {
			content := strings.Join(currentChunk, "\n")
			chunks = append(chunks, Chunk{
				Owner:    owner,
				Repo:     repo,
				Filename: filename,
				Content:  content,
				Size:     len(content),
			})

			if c.OverlapLines > 0 {
				startIdx := len(currentChunk) - c.OverlapLines
				if startIdx < 0 {
					startIdx = 0
				}
				overlapBuffer = currentChunk[startIdx:]
			}

			currentChunk = append([]string{}, overlapBuffer...)
			currentLength = 0
			for _, l := range currentChunk {
				currentLength += len(l) + 1
			}
		}

		currentChunk = append(currentChunk, line)
		currentLength += lineLen
	}

	if len(currentChunk) > 0 {
		content := strings.Join(currentChunk, "\n")
		chunks = append(chunks, Chunk{
			Owner:    owner,
			Repo:     repo,
			Filename: filename,
			Content:  content,
			Size:     len(content),
		})
	}
	return chunks
}
