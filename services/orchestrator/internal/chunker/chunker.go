package chunker

import "strings"

type Chunk struct {
	Filename string
	Content  string
	Size     int
}

type Chunker interface {
	Split(filename, text string) []Chunk
}

type CharChunker struct {
	MaxChunkSize int
	Overlap      int
}

type LineChunker struct {
	MaxChunkSize int
}

func NewCharChunker(maxChunkSize, overlap int) *CharChunker {
	return &CharChunker{
		MaxChunkSize: maxChunkSize,
		Overlap:      overlap,
	}
}

func (c *CharChunker) Split(filename, text string) []Chunk {
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
			Filename: filename,
			Content:  content,
			Size:     len(content),
		})
	}

	return chunks
}

func NewLineChunker(maxChunkSize int) *LineChunker {
	return &LineChunker{
		MaxChunkSize: maxChunkSize,
	}
}

func (c *LineChunker) Split(filename, text string) []Chunk {
	lines := strings.Split(text, "\n")
	var chunks []Chunk
	var currentChunk []string
	currentLength := 0

	for _, line := range lines {
		lineLen := len(line) + 1

		if currentLength+lineLen > c.MaxChunkSize && len(currentChunk) > 0 {
			content := strings.Join(currentChunk, "\n")
			chunks = append(chunks, Chunk{
				Filename: filename,
				Content:  content,
				Size:     len(content),
			})
			currentChunk = nil
			currentLength = 0
		}

		currentChunk = append(currentChunk, line)
		currentLength += lineLen
	}

	if len(currentChunk) > 0 {
		content := strings.Join(currentChunk, "\n")
		chunks = append(chunks, Chunk{
			Filename: filename,
			Content:  content,
			Size:     len(content),
		})
	}
	return chunks
}
