package routes

import (
	"mini_jira/contract"
	"mini_jira/handler"
	"mini_jira/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func SetupRoutes(s *contract.Service) *gin.Engine {
	r := gin.Default()

	r.RedirectTrailingSlash = false

	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	rateLimiter := mgin.NewMiddleware(instance)
	r.Use(rateLimiter)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "https://.vercel.app"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(config))

	userController := &handler.UserController{}
	userController.InitService(s)

	emailVerifController := &handler.EmailVerificationController{}
	emailVerifController.InitService(s)

	api := r.Group("/")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"status":  "success",
				"message": "Backend is running",
			})
		})
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"status":  "success",
				"message": "Api Is Healthy",
			})
		})

		// Auth routes
		api.POST("login", userController.Login)
		api.POST("register", userController.Register)
		api.POST("refresh", userController.RefreshToken)
		api.POST("logout", userController.Logout)

		// Email verification routes (public)
		api.GET("verify-email", emailVerifController.VerifyEmail)
		api.POST("resend-verification", emailVerifController.ResendVerification)

		protected := api.Group("/")
		protected.Use(middleware.Authenticate(s.Token))

		// User routes
		protected.POST("user/reset-password/:Id", userController.ResetPassword)
		protected.GET("user/:Id", userController.GetById)
		protected.GET("user/username/:username", userController.GetUserByUsername)
		protected.GET("user/email/:email", userController.GetUserByEmail)
		protected.GET("user/", userController.GetAll)
		protected.PUT("user/:Id", userController.Update)
		protected.PUT("user/:Id/status", middleware.RequireAdmin(), userController.UpdateStatus)
		protected.DELETE("user/:Id", userController.Delete)
	}
	return r
}
