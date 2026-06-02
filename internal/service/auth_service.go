package service

import (
	"context"
	"job_board/internal/domain"
	"job_board/internal/handlers/middleware"
	"job_board/internal/repository"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication related business logic.
type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, string, error) // returns a JWT token
	Refresh(ctx context.Context, oldToken string) (newAccess string, newRefresh string, err error) // returns new JWT and refresh token
	Logout(ctx context.Context, refreshToken string) error
	CleanupExpiredTokens(ctx context.Context) error
}

type authService struct {
	userRepository repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
	logger *zap.Logger
}

func NewAuthService(
	userRepository repository.UserRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	logger *zap.Logger,
	) AuthService {
	return &authService{
		userRepository: userRepository,
		refreshTokenRepository: refreshTokenRepository,
		logger: logger,
	}
}

// Register implements user registration logic
func (s *authService) Register(ctx context.Context, email, password string) error {
	// Check if user already exists
	existingUser, _ := s.userRepository.GetUserByEmail(ctx, email)

	// if user already exists, return an error
	if existingUser != nil {
		return ErrUserEmailExists
	}

	// Hash the password using a secure hashing algorithm (e.g., bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	// Create a new user object and save it to the database
	user := &domain.User{
		Email: 	  email,
		Password: string(hashedPassword),
		Role:     "user", // Default role, can be extended to support different roles 
	}

	return s.userRepository.CreateUser(ctx, user)
}

// Login implements user login logic
func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	
	// Fetch user by email
	user, err := s.userRepository.GetUserByEmail(ctx, email)

	// If an error occurs while fetching the user, return invalid credentials error
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// If user is nil, return an error
	if user == nil {
		return "", "", ErrInvalidCredentials
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	// If password does not match, return an error
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate a JWT token for the authenticated user
	accessToken, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	// Generate a refresh token for the authenticated user
	refreshToken, err := middleware.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
		
	expires := time.Now().Add(7 * 24 * time.Hour) // Set refresh token to expire in 7 days
	
	// Store the refresh token in the database
	err = s.refreshTokenRepository.SaveToken(
		ctx,
		user.ID,
		refreshToken,
		expires,
	)


	// Return the token or an error if token generation fails
	if err != nil {
		return "", "", err
	}
		

	// Return the generated token
	return accessToken, refreshToken, nil

} 

// Refresh implements token refresh logic
func (s *authService) Refresh(ctx context.Context, oldToken string) (string, string, error) {

	// This method will validate the old refresh token,
	// generate a new JWT access token and a new refresh token,
	// store the new refresh token in the database, and return both tokens to the client.

	// 1. Hash incoming refresh token to compare with stored hash in database

	// 2. Get user ID associated with the refresh token from database
	userID, err := s.refreshTokenRepository.GetUserIDByToken(ctx, oldToken)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	// 3. If no user ID found, the refresh token is invalid
	if userID == 0 {
		return "", "", ErrInvalidRefreshToken
	}

	// 4. Delete old refresh token from database to prevent reuse
	err = s.refreshTokenRepository.DeleteToken(ctx, oldToken)
	if err != nil {
		return "", "", err 
	}

	// 5. Fetch user details from database using user ID
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// 6. Generate new JWT access token
	newAccessToken, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	// Generate a new secure random refresh token
	newRefreshToken, err := middleware.GenerateRefreshToken() 
	if err != nil {
		return "", "", err
	}


	// 8. Store new hashed refresh token in database
	err = s.refreshTokenRepository.SaveToken(
		ctx,
		user.ID,
		newRefreshToken,
		time.Now().Add(7 * 24 * time.Hour),
	)
	if err != nil {
		return "", "", err
	}

	// 9. Return new access token and new refresh token to client
	return newAccessToken, newRefreshToken, nil
}

// 
func (s *authService) Logout(ctx context.Context, refreshToken string) error {

	// Delete the provided refresh token from the database to invalidate it
	return s.refreshTokenRepository.DeleteToken(ctx, refreshToken)
}

func (s *authService) CleanupExpiredTokens(ctx context.Context) error {

	// Call store to delete expired refresh tokens from database
	return s.refreshTokenRepository.DeleteExpired(ctx)
}