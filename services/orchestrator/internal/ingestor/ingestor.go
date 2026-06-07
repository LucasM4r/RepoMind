package ingestor

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"log"
	"path/filepath"

	"github.com/LucasM4r/repomind/internal/chunker"
	"github.com/LucasM4r/repomind/internal/providers"
)

type Job struct {
	Owner    string
	Repo     string
	Provider providers.Fetcher
}

func ProcessZip(resp io.ReadCloser) error {
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

		rc, err := file.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()

		if err != nil {
			return err
		}

		text := string(content)
		if len(text) == 0 {
			continue
		}

		chunker := chunker.NewLineChunker(1500)
		chunkers := chunker.Split(file.Name, text)

		log.Printf("[INFO] %s: The original file of %d bytes generated %d chunks.\n", file.Name, len(content), len(chunkers))
	}
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

func Worker(ctx context.Context, id int, jobs <-chan Job) {
	for job := range jobs {
		log.Printf("[INFO][WORKER %d] Initializing download: %s/%s\n", id, job.Owner, job.Repo)

		resp, err := job.Provider.FetchRepoZip(ctx, job.Owner, job.Repo)
		if err != nil {
			log.Printf("[ERROR][WORKER %d] Error when trying to download %s: %v \n", id, job.Repo, err)
			continue
		}

		err = ProcessZip(resp)
		resp.Close()

		if err != nil {
			log.Printf("[ERROR][WORKER %d] Error when trying to extract %s: %v\n", id, job.Repo, err)
			continue
		}

		log.Printf("[INFO][WORKER %d] Done %s/%s\n", id, job.Owner, job.Repo)
	}
}
