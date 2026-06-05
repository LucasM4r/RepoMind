package providers

import (
	"context"
	"io"
)

type Fetcher interface {
	FetchRepoZip(ctx context.Context, owner, repo string) (io.ReadCloser, error)
}
