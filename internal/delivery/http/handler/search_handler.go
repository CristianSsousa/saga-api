package handler

import (
	"net/http"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/CristianSsousa/saga-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchUC *usecase.SearchUsecase
}

func NewSearchHandler(searchUC *usecase.SearchUsecase) *SearchHandler {
	return &SearchHandler{searchUC: searchUC}
}

var validMediaTypes = map[string]domain.MediaType{
	"movie": domain.MediaTypeMovie,
	"tv":    domain.MediaTypeTV,
	"game":  domain.MediaTypeGame,
	"book":  domain.MediaTypeBook,
	"manga": domain.MediaTypeManga,
}

// Search godoc
// @Summary      Search for media
// @Description  Search across movies, TV, games, books, and manga
// @Tags         search
// @Produce      json
// @Security     BearerAuth
// @Param        q     query  string  true  "Search query"
// @Param        type  query  string  true  "Media type: movie | tv | game | book | manga"
// @Success      200   {array}   domain.Media
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	typeParam := c.Query("type")
	mt, ok := validMediaTypes[typeParam]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type; must be one of: movie, tv, game, book, manga"})
		return
	}

	results, err := h.searchUC.Search(c.Request.Context(), q, mt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	if results == nil {
		results = []domain.Media{}
	}
	c.JSON(http.StatusOK, results)
}
