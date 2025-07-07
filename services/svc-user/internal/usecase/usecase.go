package usecase

import (
	"context"
	"errors"
	"fmt"
	"ops-monorepo/services/svc-user/config"
	"ops-monorepo/services/svc-user/internal/model"
	"ops-monorepo/services/svc-user/internal/repository"
	"ops-monorepo/shared-libs/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*model.User, error)
}

type userUsecase struct {
	log  logger.Logger
	repo repository.IUserSQLRepository
	cfg  *config.Config
}

func NewUserUsecase(log logger.Logger, repo repository.IUserSQLRepository, cfg *config.Config) IUserUsecase {
	return &userUsecase{
		log:  log,
		repo: repo,
		cfg:  cfg,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// Check if user already exists
	_, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Roles:     []string{"user"},
		IsActive:  true,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Don't return password in response
	user.Password = ""

	u.log.Infof("user registered successfully: %s", user.Email)
	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// Get user by email
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := u.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Don't return password in response
	user.Password = ""

	u.log.Infof("user logged in successfully: %s", user.Email)
	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (u *userUsecase) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	// Validate refresh token
	refreshToken, err := u.repo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := u.repo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Revoke old refresh token
	if err := u.repo.RevokeRefreshToken(ctx, req.RefreshToken); err != nil {
		u.log.Errorf("failed to revoke refresh token: %v", err)
	}

	// Generate new access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := u.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	u.log.Infof("token refreshed successfully for user: %s", user.Email)
	return &model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *userUsecase) ValidateToken(ctx context.Context, tokenString string) (*model.User, error) {
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		user, err := u.repo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		// Don't return password
		user.Password = ""

		return user, nil
	}

	return nil, errors.New("invalid token")
}

func (u *userUsecase) generateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(time.Duration(u.cfg.JWT.AccessTokenDuration) * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.JWT.Secret))
}

func (u *userUsecase) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	refreshToken := &model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Duration(u.cfg.JWT.RefreshTokenDuration) * time.Hour),
		IsRevoked: false,
	}

	if err := u.repo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", err
	}

	return refreshToken.Token, nil
}