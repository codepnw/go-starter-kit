package userservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/codepnw/go-starter-kit/internal/config"
	"github.com/codepnw/go-starter-kit/internal/errs"
	"github.com/codepnw/go-starter-kit/internal/features/user"
	userrepository "github.com/codepnw/go-starter-kit/internal/features/user/repository"
	"github.com/codepnw/go-starter-kit/pkg/database"
	jwttoken "github.com/codepnw/go-starter-kit/pkg/jwttoken"
	"github.com/codepnw/go-starter-kit/pkg/utils/password"
)

type UserService interface {
	Register(ctx context.Context, u *user.User) (*UserTokenResponse, error)
	Login(ctx context.Context, email, password string) (*UserTokenResponse, error)
}

type userService struct {
	tx    database.TxManager
	token jwttoken.JWTToken
	repo  userrepository.UserRepository
}

func NewUserService(tx database.TxManager, token jwttoken.JWTToken, repo userrepository.UserRepository) UserService {
	return &userService{
		tx:    tx,
		token: token,
		repo:  repo,
	}
}

type UserTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *userService) Register(ctx context.Context, u *user.User) (*UserTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	// Check Email Exists
	exists, err := s.repo.CheckEmailExists(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.ErrEmailAlreadyExists
	}

	// Hash Password
	hashedPassword, err := password.GenerateHashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashedPassword

	var response *UserTokenResponse
	// DB Transaction
	err = s.tx.WithTx(ctx, func(tx *sql.Tx) error {
		// Insert User
		if err := s.repo.InsertUserTx(ctx, tx, u); err != nil {
			return err
		}

		// Generate Token
		resp, err := s.generateToken(u)
		if err != nil {
			return err
		}

		// Save Refresh Token
		err = s.repo.InsertRefreshTokenTx(ctx, tx, &user.RefreshToken{
			UserID:    u.ID,
			Token:     resp.RefreshToken,
			ExpiresAt: time.Now().Add(config.RefreshTokenDuration),
			Revoked:   false,
		})
		if err != nil {
			return err
		}

		response = resp
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *userService) Login(ctx context.Context, email string, pwd string) (*UserTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	// Find User Email
	foundUser, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errs.ErrInvalidEmailOrPassword
	}

	// Verify Password
	if ok := password.CompareHashedPassword(foundUser.Password, pwd); !ok {
		return nil, errs.ErrInvalidEmailOrPassword
	}

	var response *UserTokenResponse
	// DB Transaction
	err = s.tx.WithTx(ctx, func(tx *sql.Tx) error {
		// Generate Token
		resp, err := s.generateToken(foundUser)
		if err != nil {
			return err
		}

		// Save Refresh Token
		err = s.repo.InsertRefreshTokenTx(ctx, tx, &user.RefreshToken{
			UserID:    foundUser.ID,
			Token:     resp.RefreshToken,
			ExpiresAt: time.Now().Add(config.RefreshTokenDuration),
			Revoked:   false,
		})
		if err != nil {
			return err
		}

		response = resp
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *userService) generateToken(u *user.User) (*UserTokenResponse, error) {
	accessToken, err := s.token.GenerateAccessToken(u)
	if err != nil {
		return nil, fmt.Errorf("failed gen access token: %w", err)
	}

	refreshToken, err := s.token.GenerateRefreshToken(u)
	if err != nil {
		return nil, fmt.Errorf("failed gen refresh token: %w", err)
	}

	response := &UserTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}
