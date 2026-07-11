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

type DashboardHandler struct {
	activities domain.ActivityRepository
}

func NewDashboardHandler(activities domain.ActivityRepository) *DashboardHandler {
	return &DashboardHandler{activities: activities}
}

// Summary returns current-week stats (same aggregation engine as recap, but
// live/unsaved — lets the dashboard show "so far this week" without requiring
// the user to generate a recap first).
func (h *DashboardHandler) Summary(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)

	now := time.Now()
	weekStart, weekEnd := usecase.WeekBounds(now, usecase.AppTimezone)
	prevStart, prevEnd := usecase.PrevWeekBounds(weekStart, weekEnd)

	current, err := h.activities.ListByUserInRange(userID, weekStart, weekEnd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil aktivitas minggu ini")
		return
	}
	previous, err := h.activities.ListByUserInRange(userID, prevStart, prevEnd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil aktivitas minggu sebelumnya")
		return
	}

	stats := usecase.AggregateWeek(weekStart, weekEnd, usecase.AppTimezone, current, previous)
	response.OK(c, http.StatusOK, "berhasil mengambil ringkasan dashboard", stats)
}
