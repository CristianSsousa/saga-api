package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/CristianSsousa/saga-api/internal/domain"
)

type RAWGClient struct {
	apiKey string
	http   *http.Client
}

func NewRAWGClient(apiKey string) *RAWGClient {
	return &RAWGClient{apiKey: apiKey, http: &http.Client{}}
}

func (c *RAWGClient) Search(ctx context.Context, query string, _ domain.MediaType) ([]domain.Media, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.rawg.io/api/games", nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("key", c.apiKey)
	q.Set("search", query)
	q.Set("page_size", "10")
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			BackgroundImage string `json:"background_image"`
			Released        string `json:"released"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var media []domain.Media
	for _, r := range result.Results {
		media = append(media, domain.Media{
			ExternalID: fmt.Sprintf("%d", r.ID),
			Source:     "rawg",
			Type:       domain.MediaTypeGame,
			Title:      r.Name,
			CoverURL:   r.BackgroundImage,
			Metadata:   map[string]any{"released": r.Released},
		})
	}
	return media, nil
}
