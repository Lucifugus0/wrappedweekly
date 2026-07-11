package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"wrappedweekly/backend/internal/config"
	"wrappedweekly/backend/internal/middleware"
	"wrappedweekly/backend/internal/usecase"
)

type Handlers struct {
	Auth      *AuthHandler
	Activity  *ActivityHandler
	Recap     *RecapHandler
	Dashboard *DashboardHandler
}

func NewRouter(cfg config.Config, jwtManager *usecase.JWTManager, h Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendBaseURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"status": "ok"}, "message": "healthy"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/logout", h.Auth.Logout)

		// Public recap share endpoint — no auth required.
		v1.GET("/recaps/public/:slug", h.Recap.GetPublic)

		protected := v1.Group("")
		protected.Use(middleware.RequireAuth(jwtManager))
		{
			protected.GET("/auth/me", h.Auth.Me)

			activities := protected.Group("/activities")
			activities.POST("", h.Activity.Create)
			activities.GET("", h.Activity.List)
			activities.GET("/:id", h.Activity.Get)
			activities.PUT("/:id", h.Activity.Update)
			activities.DELETE("/:id", h.Activity.Delete)

			recaps := protected.Group("/recaps")
			recaps.POST("/generate", h.Recap.Generate)
			recaps.GET("", h.Recap.List)
			recaps.GET("/:id", h.Recap.Get)

			protected.GET("/dashboard/summary", h.Dashboard.Summary)
		}
	}

	return r
}
