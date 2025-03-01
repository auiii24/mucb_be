package auth

import (
	"encoding/json"
	"mucb_be/internal/domain/admin"
	"mucb_be/internal/domain/auth"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
	"strings"
	"time"

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
	isBlocked, err := u.otpAttemptRepo.CheckOtpRateLimit(req.PhoneNumber)
	if isBlocked {
		return nil, errors.NewCustomError(
			http.StatusTooManyRequests,
			"UCE002003006",
			"Too many OTP requests. Please wait before trying again.",
			err.Error(),
		)
	}

	existUser, err := u.userRepo.FindUserByPhoneNumber(req.PhoneNumber)
	if existUser != nil {
		_, err = u.sendOtpToUser(existUser)
		if err != nil {
			return nil, err
		}

		_ = u.otpAttemptRepo.IncrementOtpAttempt(req.PhoneNumber)

		phoneNumberEncrypted, err := u.encryptPhoneNumber(existUser)
		if err != nil {
			return nil, err
		}

		return &SignInUserOutput{
			Token: phoneNumberEncrypted,
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
		req.PhoneNumber,
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

	_, err = u.sendOtpToUser(user)
	if err != nil {
		return nil, err
	}

	_ = u.otpAttemptRepo.IncrementOtpAttempt(req.PhoneNumber)

	phoneNumberEncrypted, err := u.encryptPhoneNumber(user)
	if err != nil {
		return nil, err
	}

	return &SignInUserOutput{
		Token: phoneNumberEncrypted,
	}, nil
}

func (u *AuthUseCaseImpl) encryptPhoneNumber(user *user.User) (string, error) {
	phoneNumberEncrypted, err := u.encryptionService.EncryptRefreshToken(string(user.PhoneNumber))
	if err != nil {
		return "", errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005001",
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
			"UCE002002001",
			"Invalid token.",
			err.Error(),
		)
	}

	return phoneNumber, nil
}

func (u *AuthUseCaseImpl) sendOtpToUser(user *user.User) (*auth.Otp, error) {
	refCode := security.GenerateRefCode()
	otpCode := security.GenerateOtpCode()
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

func (u *AuthUseCaseImpl) VerifyOtp(req *VerifyOtpRequest) (*VerifyOtpOutput, error) {
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

	accessToken, err := u.jwtService.GenerateAccessToken(currentUser.ID.Hex(), user.RoleUser)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002004005",
			"Internal server error.",
			err.Error(),
		)
	}

	token := auth.NewToken(
		currentUser.ID,
		"",
		"TODO: USE OS NAME, DEVICE INFO",
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

func (u *AuthUseCaseImpl) RenewUser(req *RenewUserRequest, userAgent string) (*RenewUserOutput, error) {
	jsonData, err := u.encryptionService.DecryptRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005001",
			"Invalid token.",
			err.Error(),
		)
	}

	var payload RefreshTokenPayload
	err = json.Unmarshal([]byte(jsonData), &payload)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005002",
			"Invalid token.",
			err.Error(),
		)
	}

	tokenResponse, err := u.authRepo.FindTokenById(payload.Token)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005003",
			"Token not found.",
			err.Error(),
		)
	}

	if tokenResponse.UserType != auth.TokenUserType {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005004",
			"Invalid token.",
			"",
		)
	}

	existUser, err := u.userRepo.FindUserById(tokenResponse.User.Hex())
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005005",
			"User not found.",
			err.Error(),
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		existUser.ID.Hex(),
		user.RoleUser,
	)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE002005006",
			"Internal server error.",
			err.Error(),
		)
	}

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
