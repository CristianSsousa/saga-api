package domain

type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeTV    MediaType = "tv"
	MediaTypeGame  MediaType = "game"
	MediaTypeBook  MediaType = "book"
	MediaTypeManga MediaType = "manga"
)

type Media struct {
	ExternalID string         `json:"id"`
	Source     string         `json:"source"`
	Type       MediaType      `json:"type"`
	Title      string         `json:"title"`
	CoverURL   string         `json:"cover_url"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}
