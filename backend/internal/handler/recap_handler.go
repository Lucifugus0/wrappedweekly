package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wrappedweekly/backend/internal/domain"
	"wrappedweekly/backend/internal/middleware"
	"wrappedweekly/backend/internal/usecase"
	"wrappedweekly/backend/pkg/response"
)

type RecapHandler struct {
	recaps *usecase.RecapUsecase
}

func NewRecapHandler(recaps *usecase.RecapUsecase) *RecapHandler {
	return &RecapHandler{recaps: recaps}
}

type generateRecapRequest struct {
	// WeekOf is any RFC3339 date/time within the target week. Optional — defaults to now.
	WeekOf string `json:"week_of"`
	Force  bool   `json:"force"`
}

func toRecapResponse(r *domain.Recap) gin.H {
	return gin.H{
		"id":         r.ID,
		"slug":       r.Slug,
		"week_start": r.WeekStart,
		"week_end":   r.WeekEnd,
		"stats":      r.Stats,
		"narrative":  r.Narrative,
		"created_at": r.CreatedAt,
	}
}

func (h *RecapHandler) Generate(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)

	var req generateRecapRequest
	_ = c.ShouldBindJSON(&req)

	refDate := time.Now()
	if req.WeekOf != "" {
		parsed, err := time.Parse(time.RFC3339, req.WeekOf)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "week_of harus format RFC3339")
			return
		}
		refDate = parsed
	}

	recap, err := h.recaps.GenerateForWeek(userID, refDate, req.Force)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, "recap berhasil dibuat", toRecapResponse(recap))
}

func (h *RecapHandler) List(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)

	recaps, err := h.recaps.List(userID)
	if err != nil {
		handleAppError(c, err)
		return
	}

	items := make([]gin.H, 0, len(recaps))
	for _, r := range recaps {
		items = append(items, toRecapResponse(r))
	}
	response.OK(c, http.StatusOK, "berhasil mengambil daftar recap", items)
}

func (h *RecapHandler) Get(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)
	id := c.Param("id")

	recap, err := h.recaps.Get(userID, id)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "berhasil mengambil recap", toRecapResponse(recap))
}

// GetPublic is unauthenticated — mounted outside the auth-protected group.
func (h *RecapHandler) GetPublic(c *gin.Context) {
	slug := c.Param("slug")

	recap, err := h.recaps.GetPublicBySlug(slug)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "berhasil mengambil recap publik", toRecapResponse(recap))
}
