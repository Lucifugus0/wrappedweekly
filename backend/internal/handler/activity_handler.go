package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"wrappedweekly/backend/internal/domain"
	"wrappedweekly/backend/internal/middleware"
	"wrappedweekly/backend/internal/usecase"
	"wrappedweekly/backend/pkg/response"
)

type ActivityHandler struct {
	activities *usecase.ActivityUsecase
}

func NewActivityHandler(activities *usecase.ActivityUsecase) *ActivityHandler {
	return &ActivityHandler{activities: activities}
}

type activityRequest struct {
	Category   string  `json:"category" binding:"required"`
	Value      float64 `json:"value"`
	Note       *string `json:"note"`
	OccurredAt string  `json:"occurred_at" binding:"required"`
}

func toActivityResponse(a *domain.Activity) gin.H {
	return gin.H{
		"id":          a.ID,
		"category":    a.Category,
		"value":       a.Value,
		"note":        a.Note,
		"occurred_at": a.OccurredAt,
		"created_at":  a.CreatedAt,
		"updated_at":  a.UpdatedAt,
	}
}

func parseActivityRequest(req activityRequest) (usecase.ActivityInput, error) {
	occurredAt, err := time.Parse(time.RFC3339, req.OccurredAt)
	if err != nil {
		return usecase.ActivityInput{}, err
	}
	return usecase.ActivityInput{
		Category:   req.Category,
		Value:      req.Value,
		Note:       req.Note,
		OccurredAt: occurredAt,
	}, nil
}

func (h *ActivityHandler) Create(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)

	var req activityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "body tidak valid: "+err.Error())
		return
	}
	input, err := parseActivityRequest(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "occurred_at harus format RFC3339, mis. 2026-07-06T10:00:00Z")
		return
	}

	activity, err := h.activities.Create(userID, input)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusCreated, "aktivitas berhasil dicatat", toActivityResponse(activity))
}

func (h *ActivityHandler) Update(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)
	id := c.Param("id")

	var req activityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "body tidak valid: "+err.Error())
		return
	}
	input, err := parseActivityRequest(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "occurred_at harus format RFC3339, mis. 2026-07-06T10:00:00Z")
		return
	}

	activity, err := h.activities.Update(userID, id, input)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "aktivitas berhasil diperbarui", toActivityResponse(activity))
}

func (h *ActivityHandler) Delete(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)
	id := c.Param("id")

	if err := h.activities.Delete(userID, id); err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "aktivitas berhasil dihapus", nil)
}

func (h *ActivityHandler) Get(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)
	id := c.Param("id")

	activity, err := h.activities.Get(userID, id)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "berhasil mengambil aktivitas", toActivityResponse(activity))
}

func (h *ActivityHandler) List(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	filter := domain.ActivityFilter{
		Category: domain.ActivityCategory(c.Query("category")),
	}
	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filter.From = from
		} else {
			response.Error(c, http.StatusBadRequest, "from harus format RFC3339")
			return
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filter.To = to
		} else {
			response.Error(c, http.StatusBadRequest, "to harus format RFC3339")
			return
		}
	}

	activities, total, err := h.activities.List(userID, page, size, filter)
	if err != nil {
		handleAppError(c, err)
		return
	}

	items := make([]gin.H, 0, len(activities))
	for _, a := range activities {
		items = append(items, toActivityResponse(a))
	}

	response.OK(c, http.StatusOK, "berhasil mengambil daftar aktivitas", gin.H{
		"items": items,
		"page":  page,
		"size":  size,
		"total": total,
	})
}
