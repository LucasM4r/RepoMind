package domain

type RetrievedChunk struct {
	ID       string
	Owner    string
	Repo     string
	Content  string
	Distance float32
}
