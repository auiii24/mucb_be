package app

import (
	"mucb_be/internal/config"
	"mucb_be/internal/database"
	v1 "mucb_be/internal/delivery/http/v1"
	adminRepository "mucb_be/internal/infrastructure/repository/admin"
	authRepository "mucb_be/internal/infrastructure/repository/auth"
	cardRepository "mucb_be/internal/infrastructure/repository/card"
	healthScoreRepository "mucb_be/internal/infrastructure/repository/health_score"
	imageRepository "mucb_be/internal/infrastructure/repository/image"
	questionRepository "mucb_be/internal/infrastructure/repository/question"
	recordRepository "mucb_be/internal/infrastructure/repository/record"
	userRepository "mucb_be/internal/infrastructure/repository/user"
	"mucb_be/internal/infrastructure/security"
	adminUseCase "mucb_be/internal/usecase/admin"
	authUseCase "mucb_be/internal/usecase/auth"
	cardUseCase "mucb_be/internal/usecase/card"
	healthScoreUseCase "mucb_be/internal/usecase/health_score"
	imageUseCase "mucb_be/internal/usecase/image"
	questionUseCase "mucb_be/internal/usecase/question"
	recordUseCase "mucb_be/internal/usecase/record"
	userUseCase "mucb_be/internal/usecase/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type Dependencies struct {
	DBClient *mongo.Client

	JwtService        security.JwtServiceInterface
	HashService       security.HashServiceInterface
	EncryptionService security.EncryptionServiceInterface

	AdminHandlerV1       *v1.AdminHandler
	AuthHandlerV1        *v1.AuthHandler
	UserHandlerV1        *v1.UserHandler
	QuestionHandlerV1    *v1.QuestionHandler
	RecordHandlerV1      *v1.RecordHandler
	ImageHandlerV1       *v1.ImageHandler
	CardHandlerV1        *v1.CardHandler
	HealthScoreHandlerV1 *v1.HealthScoreHandler
}

func NewDependencies(cfg *config.Config, dbClient *mongo.Client) *Dependencies {

	jwtService := security.NewJwtService(cfg)
	encryptionService := security.NewEncryptionService(cfg)
	hashService := security.NewHashService()

	db := dbClient.Database(cfg.DatabaseName)
	adminCollection := db.Collection(database.AdminsCollection)
	tokenCollection := db.Collection(database.TokensCollection)
	userCollection := db.Collection(database.UsersCollection)
	otpCollection := db.Collection(database.OtpsCollection)
	otpAttemptCollection := db.Collection(database.OtpAttemptsCollection)
	questionGroupCollection := db.Collection(database.QuestionGroupsCollection)
	questionChoiceCollection := db.Collection(database.QuestionChoicesCollection)
	groupRecordCollection := db.Collection(database.GroupRecordsCollection)
	imageCollection := db.Collection(database.ImageCollection)
	cardCollection := db.Collection(database.CardCollection)
	cardRecordCollection := db.Collection(database.CardRecordsCollection)
	storyRecordCollection := db.Collection(database.StoryRecordsCollection)
	healthScoreCollection := db.Collection(database.HealthScoresCollection)

	adminRepo := adminRepository.NewAdminRepositoryMongo(adminCollection)
	authRepo := authRepository.NewAuthRepositoryMongo(tokenCollection)
	userRepo := userRepository.NewUserRepositoryMongo(userCollection)
	otpRepo := authRepository.NewOtpRepositoryMongo(otpCollection)
	otpAttemptRepo := authRepository.NewOtpAttemptRepositoryMongo(otpAttemptCollection)
	questionGroupRepo := questionRepository.NewQuestionGroupRepositoryMongo(questionGroupCollection)
	questionChoiceRepo := questionRepository.NewQuestionChoiceRepositoryMongo(questionChoiceCollection)
	groupRecordRepo := recordRepository.NewGroupRecordRepositoryMongo(groupRecordCollection)
	imageRepo := imageRepository.NewImageRepositoryMongo(imageCollection)
	cardRepo := cardRepository.NewCardRepositoryMongo(cardCollection)
	cardRecordRepo := recordRepository.NewCardRecordRepositoryMongo(cardRecordCollection)
	storyRecordRepo := recordRepository.NewStoryRecordRepositoryMongo(storyRecordCollection)
	healthScoreRepo := healthScoreRepository.NewHealthScoreRepositoryMongo(healthScoreCollection)

	adminUseCase := adminUseCase.NewAdminUseCase(adminRepo, hashService)
	authUseCase := authUseCase.NewAuthUseCase(
		userRepo,
		adminRepo,
		authRepo,
		otpRepo,
		otpAttemptRepo,
		jwtService,
		hashService,
		encryptionService,
	)
	userUseCase := userUseCase.NewUserUseCase(userRepo, groupRecordRepo, cardRecordRepo, storyRecordRepo, authRepo, jwtService)
	questionUseCase := questionUseCase.NewAdminUseCase(questionGroupRepo, questionChoiceRepo, groupRecordRepo)
	recordUseCase := recordUseCase.NewRecordUseCase(groupRecordRepo, cardRecordRepo, storyRecordRepo)
	imageUseCase := imageUseCase.NewImageUseCase(imageRepo)
	cardUseCase := cardUseCase.NewCardUseCase(cardRepo, imageRepo, cardRecordRepo)
	healthScoreUseCase := healthScoreUseCase.NewHealthScoreUseCase(healthScoreRepo, imageRepo)

	adminHandlerV1 := v1.NewAdminHandler(adminUseCase)
	authHandlerV1 := v1.NewAuthHandler(authUseCase)
	userHandlerV1 := v1.NewUserHandler(userUseCase)
	questionHandlerV1 := v1.NewQuestionHandler(questionUseCase)
	recordHandlerV1 := v1.NewRecordHandler(recordUseCase)
	imageHandlerV1 := v1.NewImageHandler(imageUseCase)
	cardHandlerV1 := v1.NewCardHandler(cardUseCase)
	healthScoreHandlerV1 := v1.NewHealthScoreHandler(healthScoreUseCase)

	return &Dependencies{
		DBClient: dbClient,

		JwtService:        jwtService,
		HashService:       hashService,
		EncryptionService: encryptionService,

		AdminHandlerV1:       adminHandlerV1,
		AuthHandlerV1:        authHandlerV1,
		UserHandlerV1:        userHandlerV1,
		QuestionHandlerV1:    questionHandlerV1,
		RecordHandlerV1:      recordHandlerV1,
		ImageHandlerV1:       imageHandlerV1,
		CardHandlerV1:        cardHandlerV1,
		HealthScoreHandlerV1: healthScoreHandlerV1,
	}
}
