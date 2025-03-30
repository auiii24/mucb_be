package auth

import (
	"encoding/json"
	"mucb_be/internal/domain/admin"
	"mucb_be/internal/domain/auth"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthUseCaseImpl struct {
	userRepo          user.UserRepository
	adminRepo         admin.AdminRepository
	authRepo          auth.AuthRepository
	otpRepo           auth.OtpRepository
	otpAttemptRepo    auth.OtpAttemptRepository
	jwtService        security.JwtServiceInterface
	hashService       security.HashServiceInterface
	encryptionService security.EncryptionServiceInterface
}

func NewAuthUseCase(
	userRepo user.UserRepository,
	adminRepo admin.AdminRepository,
	authRepo auth.AuthRepository,
	otpRepo auth.OtpRepository,
	otpAttemptRepo auth.OtpAttemptRepository,
	jwtService security.JwtServiceInterface,
	hashService security.HashServiceInterface,
	encryptionService security.EncryptionServiceInterface,
) AuthUseCaseInterface {
	return &AuthUseCaseImpl{
		userRepo:          userRepo,
		adminRepo:         adminRepo,
		authRepo:          authRepo,
		otpRepo:           otpRepo,
		otpAttemptRepo:    otpAttemptRepo,
		jwtService:        jwtService,
		hashService:       hashService,
		encryptionService: encryptionService,
	}
}

func (u *AuthUseCaseImpl) SignInAdmin(req *SignInAdminRequest, userAgent string) (*SignInAdminOutput, error) {
	admin, err := u.adminRepo.FindAdminByEmail(strings.ToLower(req.Email))
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002001001",
			"User or password incorrect.",
			err.Error(),
		)
	}

	isMatch := u.hashService.CheckHashPassword(req.Password, admin.Password)
	if !isMatch {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002001002",
			"User or password incorrect.",
			"",
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		admin.ID.Hex(),
		admin.Role,
	)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002001003",
			"User or password incorrect.",
			err.Error(),
		)
	}

	token := auth.NewToken(
		admin.ID,
		"",
		userAgent,
		auth.TokenAdminType,
	)
	err = u.authRepo.CreateToken(*token)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002001006",
			"User or password incorrect.",
			err.Error(),
		)
	}

	refreshTokenPayload := RefreshTokenPayload{
		User:  admin.ID.Hex(),
		Token: token.ID.Hex(),
		Role:  admin.Role,
	}

	refreshTokenJson, err := json.Marshal(refreshTokenPayload)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002001004",
			"User or password incorrect.",
			err.Error(),
		)
	}

	refreshToken, err := u.encryptionService.EncryptRefreshToken(string(refreshTokenJson))
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002001005",
			"User or password incorrect.",
			err.Error(),
		)
	}

	return &SignInAdminOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUseCaseImpl) RenewAdmin(req *RenewAdminRequest, userAgent string) (*RenewAdminOutput, error) {
	jsonData, err := u.encryptionService.DecryptRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002002001",
			"Invalid token.",
			err.Error(),
		)
	}

	var payload RefreshTokenPayload
	err = json.Unmarshal([]byte(jsonData), &payload)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002002002",
			"Invalid token.",
			err.Error(),
		)
	}

	tokenResponse, err := u.authRepo.FindTokenById(payload.Token)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002002003",
			"Token not found.",
			err.Error(),
		)
	}

	if tokenResponse.UserType != auth.TokenAdminType || tokenResponse.Info != userAgent {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002002004",
			"Invalid token.",
			"",
		)
	}

	adminResponse, err := u.adminRepo.FindAdminById(tokenResponse.User.Hex())
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002002005",
			"User not found.",
			err.Error(),
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		adminResponse.ID.Hex(),
		adminResponse.Role,
	)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002002006",
			"User or password incorrect.",
			err.Error(),
		)
	}

	return &RenewAdminOutput{
		AccessToken: accessToken,
	}, nil
}

func (u *AuthUseCaseImpl) SignInUser(req *SignInUserRequest) (*SignInUserOutput, error) {
	phoneNumber, err := u.encryptionService.DecryptData(req.PhoneNumber)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003007",
			"Failed to decrypt phone number",
			err.Error(),
		)
	}

	isMatched, err := user.ValidateThaiPhoneNumber(phoneNumber)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003008",
			"Failed to validate phone number",
			err.Error(),
		)
	}

	if !isMatched {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003009",
			"Invalid phone number format: "+phoneNumber+". Please enter a valid Thai phone number (must start with +66 followed by 9 digits, e.g., +66871234567).",
			"Invalid phone number format: "+phoneNumber+". Please enter a valid Thai phone number (must start with +66 followed by 9 digits, e.g., +66871234567).",
		)
	}

	var isAllowedForTester = phoneNumber == "+66800000000"

	isBlocked, err := u.otpAttemptRepo.CheckOtpRateLimit(phoneNumber)
	if isBlocked {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003006",
			"Too many OTP requests. Please wait before trying again.",
			err.Error(),
		)
	}

	existUser, err := u.userRepo.FindUserByPhoneNumber(phoneNumber)
	if existUser != nil {
		if existUser.State == user.UserStateSuspended {
			return nil, errors.NewCustomError(
				http.StatusBadRequest,
				"UCE002003010",
				"User suspended",
				"User suspended",
			)
		}

		otp, err := u.sendOtpToUser(existUser, isAllowedForTester)
		if err != nil {
			return nil, err
		}

		_ = u.otpAttemptRepo.IncrementOtpAttempt(phoneNumber)

		phoneNumberEncrypted, err := u.encryptPhoneNumber(existUser)
		if err != nil {
			return nil, err
		}

		return &SignInUserOutput{
			Token:   phoneNumberEncrypted,
			RefCode: otp.RefCode,
		}, nil
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003001",
			"Internal server error.",
			err.Error(),
		)
	}

	user := user.NewUser(
		phoneNumber,
	)

	err = u.userRepo.CreateUser(user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003002",
			"Internal server error.",
			err.Error(),
		)
	}

	otp, err := u.sendOtpToUser(user, isAllowedForTester)
	if err != nil {
		return nil, err
	}

	_ = u.otpAttemptRepo.IncrementOtpAttempt(phoneNumber)

	phoneNumberEncrypted, err := u.encryptPhoneNumber(user)
	if err != nil {
		return nil, err
	}

	return &SignInUserOutput{
		Token:   phoneNumberEncrypted,
		RefCode: otp.RefCode,
	}, nil
}

func (u *AuthUseCaseImpl) encryptPhoneNumber(user *user.User) (string, error) {
	phoneNumberEncrypted, err := u.encryptionService.EncryptRefreshToken(string(user.PhoneNumber))
	if err != nil {
		return "", errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002007001",
			"Internal server error.",
			err.Error(),
		)
	}

	return phoneNumberEncrypted, nil
}

func (u *AuthUseCaseImpl) decryptPhoneNumber(plaintext string) (string, error) {
	phoneNumber, err := u.encryptionService.DecryptRefreshToken(plaintext)
	if err != nil {
		return "", errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004012",
			"Invalid token.",
			err.Error(),
		)
	}

	return phoneNumber, nil
}

func (u *AuthUseCaseImpl) sendOtpToUser(user *user.User, isAllowedByPass bool) (*auth.Otp, error) {
	refCode := security.GenerateRefCode()

	var otpCode string
	if isAllowedByPass {
		otpCode = "111222"
	} else {
		otpCode = security.GenerateOtpCode()
	}

	otp := auth.NewOtp(
		user.ID,
		user.PhoneNumber,
		refCode,
		otpCode,
	)

	err := u.otpRepo.CreateOtp(otp)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002003005",
			"Can not send OTP.",
			err.Error(),
		)
	}

	return otp, nil
}

func (u *AuthUseCaseImpl) VerifyOtp(req *VerifyOtpRequest, c *gin.Context) (*VerifyOtpOutput, error) {
	phoneNumber, err := u.decryptPhoneNumber(req.Token)
	if err != nil {
		return nil, err
	}

	latestOtp, err := u.otpRepo.FindLatestOtpByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004001",
			"Internal server error.",
			err.Error(),
		)
	}

	if time.Since(latestOtp.CreatedAt) > 5*time.Minute {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004002",
			"OTP has expired.",
			"",
		)
	}

	if latestOtp.IsUsed {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004011",
			"This OTP has already been used.",
			"",
		)
	}

	if latestOtp.AttemptCount >= 5 {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004010",
			"Too many incorrect OTP attempts. Please request a new OTP.",
			"",
		)
	}

	if latestOtp.Code != req.Code {
		u.otpRepo.IncrementOtpAttemptsById(latestOtp.ID.Hex())

		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004003",
			"Invalid OTP code.",
			"",
		)
	}

	currentUser, err := u.userRepo.FindUserByPhoneNumber(phoneNumber)
	if latestOtp.Code != req.Code {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004004",
			"User not found.",
			err.Error(),
		)
	}

	if currentUser.State == user.UserStateSuspended {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004009",
			"User suspended.",
			"User suspended.",
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(currentUser.ID.Hex(), user.RoleUser)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004005",
			"Internal server error.",
			err.Error(),
		)
	}

	uniqueId := c.GetHeader("X-UNIQUE-ID")
	brand := c.GetHeader("X-BRAND")
	platform := c.GetHeader("X-PLATFORM")
	deviceId := c.GetHeader("X-DEVICE-ID")

	deviceInfo := platform + "\n" + brand + "\n" + deviceId

	token := auth.NewToken(
		currentUser.ID,
		uniqueId,
		deviceInfo,
		auth.TokenUserType,
	)
	err = u.authRepo.CreateToken(*token)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004006",
			"Internal server error.",
			err.Error(),
		)
	}

	refreshTokenPayload := RefreshTokenPayload{
		User:  currentUser.ID.Hex(),
		Token: token.ID.Hex(),
		Role:  user.RoleUser,
	}

	refreshTokenJson, err := json.Marshal(refreshTokenPayload)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004007",
			"Internal server error.",
			err.Error(),
		)
	}

	refreshToken, err := u.encryptionService.EncryptRefreshToken(string(refreshTokenJson))
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004008",
			"Internal server error.",
			err.Error(),
		)
	}

	_ = u.otpRepo.MarkOtpAsUsedById(latestOtp.ID.Hex())
	_ = u.otpAttemptRepo.ClearOtpAttempts(phoneNumber)

	return &VerifyOtpOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Status:       currentUser.State,
	}, nil
}

func (u *AuthUseCaseImpl) RenewUser(req *RenewUserRequest, c *gin.Context) (*RenewUserOutput, error) {
	jsonData, err := u.encryptionService.DecryptRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005001",
			"Invalid token.",
			err.Error(),
		)
	}

	var payload RefreshTokenPayload
	err = json.Unmarshal([]byte(jsonData), &payload)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005002",
			"Invalid token.",
			err.Error(),
		)
	}

	tokenResponse, err := u.authRepo.FindTokenById(payload.Token)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005003",
			"Token not found.",
			err.Error(),
		)
	}

	if tokenResponse.UserType != auth.TokenUserType {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005004",
			"Invalid token.",
			"",
		)
	}

	existUser, err := u.userRepo.FindUserById(tokenResponse.User.Hex())
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005005",
			"User not found.",
			err.Error(),
		)
	}

	if existUser.State == user.UserStateSuspended {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005007",
			"User not found.",
			"User not found.",
		)
	}

	uniqueId := c.GetHeader("X-UNIQUE-ID")

	if tokenResponse.UniqueID != uniqueId {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005008",
			"Invalid info.",
			"Invalid info.",
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		existUser.ID.Hex(),
		user.RoleUser,
	)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusUnauthorized,
			"UCE002005006",
			"Internal server error.",
			err.Error(),
		)
	}

	u.authRepo.UpdateTimestampByTokenId(payload.Token)

	return &RenewUserOutput{
		AccessToken: accessToken,
	}, nil
}

func (u *AuthUseCaseImpl) SignOut(req *SignOutRequest) error {
	jsonData, err := u.encryptionService.DecryptRefreshToken(req.RefreshToken)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002006001",
			"Invalid token.",
			err.Error(),
		)
	}

	var payload RefreshTokenPayload
	err = json.Unmarshal([]byte(jsonData), &payload)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002006002",
			"Invalid token.",
			err.Error(),
		)
	}

	u.authRepo.RemoveTokenById(payload.Token)

	return nil
}

func (u *AuthUseCaseImpl) FindAllToken(req *FindAllTokensRequest, claims *security.AccessTokenModel) (*FindAllTokensOutput, error) {
	jsonData, err := u.encryptionService.DecryptRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002008001",
			"Invalid token.",
			err.Error(),
		)
	}

	var payload RefreshTokenPayload
	err = json.Unmarshal([]byte(jsonData), &payload)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002008002",
			"Invalid token.",
			err.Error(),
		)
	}

	tokens, err := u.authRepo.FindAllTokenByUserId(claims.ID)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002008003",
			"Failed to get token.",
			err.Error(),
		)
	}

	var bundleTokens []TokenWithFlag
	for _, token := range *tokens {
		bundleTokens = append(bundleTokens, TokenWithFlag{
			IsCurrentToken: payload.Token == string(token.ID.Hex()),
			Token:          token,
		})
	}

	// ✅ First: Sort by IsCurrentToken (true first)
	sort.SliceStable(bundleTokens, func(i, j int) bool {
		return bundleTokens[i].IsCurrentToken && !bundleTokens[j].IsCurrentToken
	})

	// ✅ Second: Sort by UpdatedAt (most recent first) within each IsCurrentToken group
	sort.SliceStable(bundleTokens, func(i, j int) bool {
		if bundleTokens[i].IsCurrentToken == bundleTokens[j].IsCurrentToken {
			return bundleTokens[i].UpdatedAt.After(bundleTokens[j].UpdatedAt) // Newest first
		}
		return false // Keep existing order from the first sort
	})

	return &FindAllTokensOutput{
		Tokens: &bundleTokens,
	}, nil
}

// RevokeToken implements AuthUseCaseInterface.
func (u *AuthUseCaseImpl) RevokeToken(req *RevokeTokenRequest, claims *security.AccessTokenModel) error {
	existToken, err := u.authRepo.FindTokenById(req.Token)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002009001",
			"Failed to get token.",
			err.Error(),
		)
	}

	if existToken.User.Hex() != claims.ID {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002009002",
			"Failed to validate user.",
			"Failed to validate user.",
		)
	}

	err = u.authRepo.RemoveTokenById(req.Token)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002009003",
			"Failed to remove token.",
			err.Error(),
		)
	}

	return nil
}
