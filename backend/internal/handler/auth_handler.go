package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"wrappedweekly/backend/internal/config"
	"wrappedweekly/backend/internal/domain"
	"wrappedweekly/backend/internal/middleware"
	"wrappedweekly/backend/internal/usecase"
	"wrappedweekly/backend/pkg/apperror"
	"wrappedweekly/backend/pkg/response"
)

type AuthHandler struct {
	auth *usecase.AuthUsecase
	cfg  config.Config
}

func NewAuthHandler(auth *usecase.AuthUsecase, cfg config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, cfg: cfg}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func toUserResponse(u *domain.User) gin.H {
	return gin.H{
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"created_at": u.CreatedAt,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "body tidak valid: "+err.Error())
		return
	}

	user, err := h.auth.Register(usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		handleAppError(c, err)
		return
	}

	response.OK(c, http.StatusCreated, "registrasi berhasil", toUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "body tidak valid: "+err.Error())
		return
	}

	result, err := h.auth.Login(usecase.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		handleAppError(c, err)
		return
	}

	h.setAuthCookie(c, result.Token)
	response.OK(c, http.StatusOK, "login berhasil", gin.H{
		"user": toUserResponse(result.User),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.UserIDFromContext(c)
	user, err := h.auth.Me(userID)
	if err != nil {
		handleAppError(c, err)
		return
	}
	response.OK(c, http.StatusOK, "berhasil mengambil profil", toUserResponse(user))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.clearAuthCookie(c)
	response.OK(c, http.StatusOK, "logout berhasil", nil)
}

func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	maxAge := int(h.cfg.JWTExpiry.Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.CookieName, token, maxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

func (h *AuthHandler) clearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.CookieName, "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

func handleAppError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.Status, appErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, "terjadi kesalahan internal")
}
