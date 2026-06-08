package ingestor

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"log"
	"path/filepath"
	"strconv"

	"github.com/LucasM4r/repomind/internal/chunker"
	"github.com/LucasM4r/repomind/internal/providers"
	"github.com/LucasM4r/repomind/internal/rpc"
	"github.com/LucasM4r/repomind/internal/vector_db"
)

type Job struct {
	Owner    string
	Repo     string
	Provider providers.Fetcher
}

type Ingestor struct {
	AIClient *rpc.AIClient
	Repo     *vector_db.VectorRepository
}

func NewIngestor(ai *rpc.AIClient, repo *vector_db.VectorRepository) *Ingestor {
	return &Ingestor{AIClient: ai, Repo: repo}
}

func Worker(ctx context.Context, id int, jobs <-chan Job, ing *Ingestor) {
	for job := range jobs {
		log.Printf("[INFO][WORKER %d] Initializing download: %s/%s\n", id, job.Owner, job.Repo)

		resp, err := job.Provider.FetchRepoZip(ctx, job.Owner, job.Repo)
		if err != nil {
			log.Printf("[ERROR][WORKER %d] Error downloading %s: %v\n", id, job.Repo, err)
			continue
		}

		if err := ing.ProcessZip(ctx, resp); err != nil {
			log.Printf("[ERROR][WORKER %d] Error processing zip %s: %v\n", id, job.Repo, err)
		}
		resp.Close()

		log.Printf("[INFO][WORKER %d] Done %s/%s\n", id, job.Owner, job.Repo)
	}
}

func (in *Ingestor) ProcessZip(ctx context.Context, resp io.ReadCloser) error {
	bodyBytes, err := io.ReadAll(resp)
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() || !isCodeFile(file.Name) {
			continue
		}

		if err := in.processFile(ctx, file); err != nil {
			log.Printf("[WARNING] Failed to process file %s: %v", file.Name, err)
		}
	}
	return nil
}

func (in *Ingestor) processFile(ctx context.Context, file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	text := string(content)
	if len(text) == 0 {
		return nil
	}

	c := chunker.NewLineChunker(1500)
	chunks := c.Split(file.Name, text)

	if len(chunks) > 0 {
		return in.embedAndSave(ctx, file.Name, chunks)
	}
	return nil
}

func (in *Ingestor) embedAndSave(ctx context.Context, filename string, chunks []chunker.Chunk) error {
	var texts []string
	for _, chunk := range chunks {
		texts = append(texts, chunk.Content)
	}

	respAI, err := in.AIClient.GetEmbeddings(ctx, texts)
	if err != nil {
		return err
	}

	for i, chunk := range chunks {
		chunkID := filename + "_chunk_" + strconv.Itoa(i)

		if i < len(respAI.Embeddings) {
			embedding := respAI.Embeddings[i].Values
			if err := in.Repo.SaveChunk(ctx, chunkID, chunk.Content, embedding); err != nil {
				log.Printf("[ERROR] Failed to save chunk %s to DB: %v", chunkID, err)
			}
		}
	}

	log.Printf("[INFO] %s: Successfully processed and saved %d chunks.\n", filename, len(chunks))
	return nil
}

func isCodeFile(name string) bool {
	ext := filepath.Ext(name)
	switch ext {
	case ".go", ".md", ".py", ".js", ".ts", ".c":
		return true
	default:
		return false
	}
}
