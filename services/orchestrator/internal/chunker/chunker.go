package chunker

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/LucasM4r/repomind/internal/domain"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/cpp"
	"github.com/smacker/go-tree-sitter/csharp"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/php"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type Chunker interface {
	Split(owner, repo, filename, text string) []domain.RawChunk
}

func createChunk(owner, repo, filename, content string) domain.RawChunk {
	return domain.RawChunk{
		Owner:    owner,
		Repo:     repo,
		Filename: filename,
		Content:  content,
		Size:     len(content),
	}
}

type CharChunker struct {
	MaxChunkSize int
	Overlap      int
}

type LineChunker struct {
	MaxChunkSize int
	OverlapLines int
}

type TreeSitterChunker struct {
	MaxChunkSize int
	parsers      map[string]*sitter.Language
	fallback     Chunker
}

func NewCharChunker(maxChunkSize, overlap int) *CharChunker {
	return &CharChunker{
		MaxChunkSize: maxChunkSize,
		Overlap:      overlap,
	}
}

func (c *CharChunker) Split(owner, repo, filename, text string) []domain.RawChunk {
	runes := []rune(text)
	var chunks []domain.RawChunk
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
		chunks = append(chunks, createChunk(owner, repo, filename, content))
	}

	return chunks
}

func NewLineChunker(maxChunkSize, overlapLines int) *LineChunker {
	return &LineChunker{
		MaxChunkSize: maxChunkSize,
		OverlapLines: overlapLines,
	}
}

func (c *LineChunker) Split(owner, repo, filename, text string) []domain.RawChunk {
	lines := strings.Split(text, "\n")
	var chunks []domain.RawChunk
	var currentChunk []string
	var overlapBuffer []string
	currentLength := 0

	for _, line := range lines {
		lineLen := len(line) + 1

		if currentLength+lineLen > c.MaxChunkSize && len(currentChunk) > 0 {
			content := strings.Join(currentChunk, "\n")
			chunks = append(chunks, createChunk(owner, repo, filename, content))

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
		chunks = append(chunks, createChunk(owner, repo, filename, content))
	}
	return chunks
}

func NewTreeSitterChunker(maxChunkSize int, fallback Chunker) *TreeSitterChunker {
	return &TreeSitterChunker{
		MaxChunkSize: maxChunkSize,
		fallback:     fallback,
		parsers: map[string]*sitter.Language{
			// Go
			".go": golang.GetLanguage(),
			// Python
			".py": python.GetLanguage(),
			// JavaScript / TypeScript
			".js":  javascript.GetLanguage(),
			".jsx": javascript.GetLanguage(),
			".ts":  typescript.GetLanguage(),
			".tsx": tsx.GetLanguage(),
			// C / C++
			".c":   c.GetLanguage(),
			".h":   c.GetLanguage(),
			".cpp": cpp.GetLanguage(),
			".hpp": cpp.GetLanguage(),
			".cc":  cpp.GetLanguage(),
			// Java
			".java": java.GetLanguage(),
			// C#
			".cs": csharp.GetLanguage(),
			// Rust
			".rs": rust.GetLanguage(),
			// Ruby
			".rb": ruby.GetLanguage(),
			// PHP
			".php": php.GetLanguage(),
			// Shell/Bash
			".sh":   bash.GetLanguage(),
			".bash": bash.GetLanguage(),
		},
	}
}

func (c *TreeSitterChunker) Split(owner, repo, filename, text string) []domain.RawChunk {
	ext := filepath.Ext(filename)
	lang, isSupported := c.parsers[ext]

	if !isSupported {
		if c.fallback != nil {
			return c.fallback.Split(owner, repo, filename, text)
		}
		return nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(text))
	if err != nil {
		if c.fallback != nil {
			return c.fallback.Split(owner, repo, filename, text)
		}
		return nil
	}

	root := tree.RootNode()
	var chunks []domain.RawChunk
	var currentContent string

	for i := 0; i < int(root.ChildCount()); i++ {
		child := root.Child(i)
		nodeContent := child.Content([]byte(text))

		if len(nodeContent) > c.MaxChunkSize {
			if len(currentContent) > 0 {
				chunks = append(chunks, createChunk(owner, repo, filename, currentContent))
				currentContent = ""
			}

			if c.fallback != nil {
				fallbackChunks := c.fallback.Split(owner, repo, filename, nodeContent)
				chunks = append(chunks, fallbackChunks...)
			} else {
				chunks = append(chunks, createChunk(owner, repo, filename, nodeContent))
			}
			continue
		}

		if len(currentContent)+len(nodeContent) > c.MaxChunkSize && len(currentContent) > 0 {
			chunks = append(chunks, createChunk(owner, repo, filename, currentContent))
			currentContent = ""
		}

		if len(currentContent) > 0 {
			currentContent += "\n\n" + nodeContent
		} else {
			currentContent = nodeContent
		}
	}

	if len(currentContent) > 0 {
		chunks = append(chunks, createChunk(owner, repo, filename, currentContent))
	}

	return chunks
}
