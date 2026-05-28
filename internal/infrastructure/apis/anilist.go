package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CristianSsousa/saga-api/internal/domain"
)

type AniListClient struct {
	http *http.Client
}

func NewAniListClient() *AniListClient {
	return &AniListClient{http: &http.Client{}}
}

const aniListQuery = `
query ($search: String) {
  Page(perPage: 10) {
    media(search: $search, type: MANGA) {
      id
      title { romaji english }
      coverImage { large }
      startDate { year }
      genres
    }
  }
}`

func (c *AniListClient) Search(ctx context.Context, query string, _ domain.MediaType) ([]domain.Media, error) {
	payload := map[string]any{
		"query":     aniListQuery,
		"variables": map[string]any{"search": query},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://graphql.anilist.co", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Page struct {
				Media []struct {
					ID    int `json:"id"`
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
					} `json:"title"`
					CoverImage struct {
						Large string `json:"large"`
					} `json:"coverImage"`
					StartDate struct {
						Year int `json:"year"`
					} `json:"startDate"`
					Genres []string `json:"genres"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var media []domain.Media
	for _, m := range result.Data.Page.Media {
		title := m.Title.English
		if title == "" {
			title = m.Title.Romaji
		}
		media = append(media, domain.Media{
			ExternalID: fmt.Sprintf("%d", m.ID),
			Source:     "anilist",
			Type:       domain.MediaTypeManga,
			Title:      title,
			CoverURL:   m.CoverImage.Large,
			Metadata:   map[string]any{"year": m.StartDate.Year, "genres": m.Genres},
		})
	}
	return media, nil
}
