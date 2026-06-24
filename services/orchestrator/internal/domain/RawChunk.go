package domain

type RawChunk struct {
	Owner    string
	Repo     string
	Filename string
	Content  string
	Size     int
}
