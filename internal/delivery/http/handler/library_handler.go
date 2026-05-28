package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/CristianSsousa/saga-api/internal/usecase"
	"github.com/CristianSsousa/saga-api/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	libraryUC *usecase.LibraryUsecase
}

func NewLibraryHandler(libraryUC *usecase.LibraryUsecase) *LibraryHandler {
	return &LibraryHandler{libraryUC: libraryUC}
}

type addLibraryRequest struct {
	Media  domain.Media  `json:"media"  binding:"required"`
	Status domain.Status `json:"status" binding:"required"`
}

type updateLibraryRequest struct {
	Status *domain.Status `json:"status"`
	Rating *int           `json:"user_rating"`
	Notes  *string        `json:"notes"`
}

// ListLibrary godoc
// @Summary      List library items
// @Description  Returns paginated library items for the authenticated user
// @Tags         library
// @Produce      json
// @Security     BearerAuth
// @Param        cursor  query  string  false  "Pagination cursor"
// @Param        limit   query  int     false  "Page size (default 20)"
// @Param        type    query  string  false  "Filter by media type"
// @Param        status  query  string  false  "Filter by status"
// @Success      200  {object}  pagination.Page
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/library [get]
func (h *LibraryHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	limit := pagination.DefaultLimit
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	filter := domain.ListFilter{
		Cursor: c.Query("cursor"),
		Limit:  limit + 1, // fetch one extra to determine HasMore
	}
	if t := c.Query("type"); t != "" {
		mt := domain.MediaType(t)
		filter.MediaType = &mt
	}
	if s := c.Query("status"); s != "" {
		st := domain.Status(s)
		filter.Status = &st
	}

	items, err := h.libraryUC.List(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list library"})
		return
	}
	if items == nil {
		items = []domain.LibraryItem{}
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	nextCursor := ""
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].ID
	}

	c.JSON(http.StatusOK, pagination.Page{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
}

// AddToLibrary godoc
// @Summary      Add media to library
// @Description  Add a media item to the authenticated user's library
// @Tags         library
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  addLibraryRequest  true  "Library item"
// @Success      201   {object}  domain.LibraryItem
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/library [post]
func (h *LibraryHandler) Add(c *gin.Context) {
	userID := c.GetString("user_id")
	var req addLibraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.libraryUC.Add(c.Request.Context(), userID, usecase.AddToLibraryInput{
		Media:  req.Media,
		Status: req.Status,
	})
	if err != nil {
		log.Printf("[library.Add] userID=%s media=%+v status=%s err=%v", userID, req.Media, req.Status, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to library"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateLibraryItem godoc
// @Summary      Update a library item
// @Description  Update status, rating, or notes of a library entry
// @Tags         library
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string               true  "Library item ID"
// @Param        body  body  updateLibraryRequest true  "Fields to update"
// @Success      200   {object}  domain.LibraryItem
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/library/{id} [put]
func (h *LibraryHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	var req updateLibraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.libraryUC.Update(c.Request.Context(), id, userID, usecase.UpdateLibraryInput{
		Status: req.Status,
		Rating: req.Rating,
		Notes:  req.Notes,
	})
	if err != nil {
		if errors.Is(err, domain.ErrLibraryItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteLibraryItem godoc
// @Summary      Remove a library item
// @Description  Soft-delete a library entry
// @Tags         library
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Library item ID"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/library/{id} [delete]
func (h *LibraryHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	err := h.libraryUC.Delete(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrLibraryItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}
	c.Status(http.StatusNoContent)
}
