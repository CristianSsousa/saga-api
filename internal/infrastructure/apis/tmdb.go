package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/CristianSsousa/saga-api/internal/domain"
)

type TMDBClient struct {
	apiKey string
	http   *http.Client
}

func NewTMDBClient(apiKey string) *TMDBClient {
	return &TMDBClient{apiKey: apiKey, http: &http.Client{}}
}

func (c *TMDBClient) Search(ctx context.Context, query string, mediaType domain.MediaType) ([]domain.Media, error) {
	endpoint := "https://api.themoviedb.org/3/search/multi"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("api_key", c.apiKey)
	q.Set("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			ID           int    `json:"id"`
			Title        string `json:"title"`
			Name         string `json:"name"`
			PosterPath   string `json:"poster_path"`
			MediaType    string `json:"media_type"`
			ReleaseDate  string `json:"release_date"`
			FirstAirDate string `json:"first_air_date"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var media []domain.Media
	for _, r := range result.Results {
		if r.MediaType != "movie" && r.MediaType != "tv" {
			continue
		}
		title := r.Title
		if title == "" {
			title = r.Name
		}
		year := r.ReleaseDate
		if year == "" {
			year = r.FirstAirDate
		}
		coverURL := ""
		if r.PosterPath != "" {
			coverURL = "https://image.tmdb.org/t/p/w500" + r.PosterPath
		}
		mt := domain.MediaTypeMovie
		if r.MediaType == "tv" {
			mt = domain.MediaTypeTV
		}
		media = append(media, domain.Media{
			ExternalID: fmt.Sprintf("%d", r.ID),
			Source:     "tmdb",
			Type:       mt,
			Title:      title,
			CoverURL:   coverURL,
			Metadata:   map[string]any{"year": year},
		})
	}
	return media, nil
}
