package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/CristianSsousa/saga-api/internal/domain"
)

type OpenLibraryClient struct {
	http *http.Client
}

func NewOpenLibraryClient() *OpenLibraryClient {
	return &OpenLibraryClient{http: &http.Client{}}
}

func (c *OpenLibraryClient) Search(ctx context.Context, query string, _ domain.MediaType) ([]domain.Media, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openlibrary.org/search.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "saga-app/1.0 (cristian.sousa365@outlook.com)")
	q := url.Values{}
	q.Set("q", query)
	q.Set("limit", "10")
	q.Set("fields", "key,title,author_name,cover_i,first_publish_year")
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Docs []struct {
			Key              string   `json:"key"`
			Title            string   `json:"title"`
			AuthorName       []string `json:"author_name"`
			CoverI           int      `json:"cover_i"`
			FirstPublishYear int      `json:"first_publish_year"`
		} `json:"docs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var media []domain.Media
	for _, d := range result.Docs {
		coverURL := ""
		if d.CoverI != 0 {
			coverURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", d.CoverI)
		}
		author := ""
		if len(d.AuthorName) > 0 {
			author = d.AuthorName[0]
		}
		media = append(media, domain.Media{
			ExternalID: d.Key,
			Source:     "openlibrary",
			Type:       domain.MediaTypeBook,
			Title:      d.Title,
			CoverURL:   coverURL,
			Metadata:   map[string]any{"author": author, "year": d.FirstPublishYear},
		})
	}
	return media, nil
}
