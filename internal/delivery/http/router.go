package http

import (
	"mucb_be/internal/app"
	"mucb_be/internal/config"
	"mucb_be/internal/domain/admin"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, cfg *config.Config, deps *app.Dependencies) {
	router.Use(middleware.BasicAuthMiddleware(cfg.ApiKey))
	router.Use(middleware.RequestLimitMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.NoRoute(middleware.InvalidEndpointMiddleware())

	allowedOnlySuperAdminRole := middleware.SpecificAuthMiddleware(deps.JwtService, []string{admin.RoleSuperAdmin})
	allowedOnlyAdminRole := middleware.SpecificAuthMiddleware(deps.JwtService, []string{admin.RoleSuperAdmin, admin.RoleAdmin})
	allowedOnlyUserRole := middleware.SpecificAuthMiddleware(deps.JwtService, []string{user.RoleUser})
	allowedAllRole := middleware.SpecificAuthMiddleware(deps.JwtService, []string{admin.RoleSuperAdmin, admin.RoleAdmin, user.RoleUser})

	api := router.Group("/api")
	routesV1 := api.Group("/v1")

	authRoutesV1 := routesV1.Group("/auth")
	authRoutesV1.POST("/sign-in-admin", deps.AuthHandlerV1.SignInAdmin)
	authRoutesV1.POST("/renew-admin", deps.AuthHandlerV1.RenewAdmin)
	authRoutesV1.POST("/sign-in", deps.AuthHandlerV1.SignInUser)
	authRoutesV1.POST("/verify-otp", deps.AuthHandlerV1.VerifyOtpUser)
	authRoutesV1.POST("/renew", deps.AuthHandlerV1.RenewUser)
	authRoutesV1.DELETE("/sign-out", allowedAllRole, deps.AuthHandlerV1.SignOut)

	adminRoutesV1 := routesV1.Group("/admin")
	adminRoutesV1.POST("/create", allowedOnlySuperAdminRole, deps.AdminHandlerV1.CreateAdmin)

	userRoutesV1 := routesV1.Group("/user")
	userRoutesV1.PUT("/update-info", allowedOnlyUserRole, deps.UserHandlerV1.UpdateUserInfo)

	questionRoutesV1 := routesV1.Group("/question")
	questionRoutesV1.POST("/create-group", allowedOnlyAdminRole, deps.QuestionHandlerV1.CreateQuestionGroup)
	questionRoutesV1.GET("/group", allowedOnlyAdminRole, deps.QuestionHandlerV1.GetAllQuestionGroups)
	questionRoutesV1.POST("/create-choice", allowedOnlyAdminRole, deps.QuestionHandlerV1.CreateQuestionChoice)
	questionRoutesV1.POST("/choice", allowedOnlyAdminRole, deps.QuestionHandlerV1.FindAllQuestionChoiceByQuestionGroup)
	questionRoutesV1.GET("/exam", allowedOnlyUserRole, deps.QuestionHandlerV1.GetQuestionWithRandomChoices)
	questionRoutesV1.PUT("/update-choice", allowedOnlyAdminRole, deps.QuestionHandlerV1.UpdateQuestion)

	recordRoutesV1 := routesV1.Group("/record")
	recordRoutesV1.POST("/submit-group-answer", allowedOnlyUserRole, deps.RecordHandlerV1.SubmitGroupAnswer)
	recordRoutesV1.POST("/submit-card-answer", allowedOnlyUserRole, deps.RecordHandlerV1.SubmitCardAnswer)
	recordRoutesV1.POST("/submit-story-answer", allowedOnlyUserRole, deps.RecordHandlerV1.SubmitStoryAnswer)

	imageRoutesV1 := routesV1.Group("/image")
	imageRoutesV1.POST("/upload", allowedOnlyAdminRole, deps.ImageHandlerV1.UploadImage)
	imageRoutesV1.GET("/:imageId", deps.ImageHandlerV1.GetImage)

	cardRoutesV1 := routesV1.Group("/card")
	cardRoutesV1.POST("/create", allowedOnlyAdminRole, deps.CardHandlerV1.CreateCard)
	cardRoutesV1.GET("/:cardId", allowedOnlyAdminRole, deps.CardHandlerV1.GetCard)
	cardRoutesV1.GET("/", allowedAllRole, deps.CardHandlerV1.GetAllCards)
	cardRoutesV1.POST("/activate", allowedOnlyAdminRole, deps.CardHandlerV1.ActivateCard)
	cardRoutesV1.PUT("/update", allowedOnlyAdminRole, deps.CardHandlerV1.UpdateCard)

	healthScoreRoutesV1 := routesV1.Group("/health-score")
	healthScoreRoutesV1.POST("/create", allowedOnlyAdminRole, deps.HealthScoreHandlerV1.CreateHealthScore)
	healthScoreRoutesV1.GET("/", allowedOnlyAdminRole, deps.HealthScoreHandlerV1.GetAllHealthScore)
	healthScoreRoutesV1.GET("/:healthScoreId", allowedOnlyAdminRole, deps.HealthScoreHandlerV1.GetHealthScore)
	healthScoreRoutesV1.PUT("/update", allowedOnlyAdminRole, deps.HealthScoreHandlerV1.UpdateHealthScoreById)
	healthScoreRoutesV1.POST("/score", allowedOnlyUserRole, deps.HealthScoreHandlerV1.GetContentByScore)
}
