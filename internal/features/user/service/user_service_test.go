package userservice_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/codepnw/go-starter-kit/internal/errs"
	"github.com/codepnw/go-starter-kit/internal/features/user"
	userrepository "github.com/codepnw/go-starter-kit/internal/features/user/repository"
	userservice "github.com/codepnw/go-starter-kit/internal/features/user/service"
	"github.com/codepnw/go-starter-kit/pkg/database"
	jwttoken "github.com/codepnw/go-starter-kit/pkg/jwttoken"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ErrDB = errors.New("DB Error")

func TestRegister(t *testing.T) {
	type testCase struct {
		name        string
		input       *user.User
		mockFn      func(tx *database.MockTxManager, token *jwttoken.MockJWTToken, repo *userrepository.MockUserRepository, input *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(tx *database.MockTxManager, token *jwttoken.MockJWTToken, repo *userrepository.MockUserRepository, input *user.User) {
				repo.EXPECT().CheckEmailExists(gomock.Any(), input.Email).Return(false, nil).Times(1)

				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				repo.EXPECT().InsertUserTx(gomock.Any(), nil, input).Return(nil).Times(1)

				token.EXPECT().GenerateAccessToken(input).Return("mock-access-token", nil).Times(1)
				token.EXPECT().GenerateRefreshToken(input).Return("mock-refresh-token", nil).Times(1)

				repo.EXPECT().InsertRefreshTokenTx(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail email exists",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(tx *database.MockTxManager, token *jwttoken.MockJWTToken, repo *userrepository.MockUserRepository, input *user.User) {
				repo.EXPECT().CheckEmailExists(gomock.Any(), input.Email).Return(true, errs.ErrEmailAlreadyExists).Times(1)
			},
			expectedErr: errs.ErrEmailAlreadyExists,
		},
		{
			name:  "fail insert user",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(tx *database.MockTxManager, token *jwttoken.MockJWTToken, repo *userrepository.MockUserRepository, input *user.User) {
				repo.EXPECT().CheckEmailExists(gomock.Any(), input.Email).Return(false, nil).Times(1)

				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				repo.EXPECT().InsertUserTx(gomock.Any(), nil, input).Return(ErrDB).Times(1)
			},
			expectedErr: ErrDB,
		},
		{
			name:  "fail generate token",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(tx *database.MockTxManager, token *jwttoken.MockJWTToken, repo *userrepository.MockUserRepository, input *user.User) {
				repo.EXPECT().CheckEmailExists(gomock.Any(), input.Email).Return(false, nil).Times(1)

				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				repo.EXPECT().InsertUserTx(gomock.Any(), nil, input).Return(nil).Times(1)

				token.EXPECT().GenerateAccessToken(input).Return("mock-access-token", nil).Times(1)
				token.EXPECT().GenerateRefreshToken(input).Return("", ErrDB).Times(1)
			},
			expectedErr: ErrDB,
		},
		{
			name:  "fail insert token",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(tx *database.MockTxManager, token *jwttoken.MockJWTToken, repo *userrepository.MockUserRepository, input *user.User) {
				repo.EXPECT().CheckEmailExists(gomock.Any(), input.Email).Return(false, nil).Times(1)

				tx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				repo.EXPECT().InsertUserTx(gomock.Any(), nil, input).Return(nil).Times(1)

				token.EXPECT().GenerateAccessToken(input).Return("mock-access-token", nil).Times(1)
				token.EXPECT().GenerateRefreshToken(input).Return("mock-refresh-token", nil).Times(1)

				repo.EXPECT().InsertRefreshTokenTx(gomock.Any(), nil, gomock.Any()).Return(ErrDB).Times(1)
			},
			expectedErr: ErrDB,
		},
	}

	for _, tc := range testCases {
		mockToken, mockTx, mockRepo, service := setup(t)

		tc.mockFn(mockTx, mockToken, mockRepo, tc.input)

		resp, err := service.Register(context.Background(), tc.input)

		if tc.expectedErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		}
	}
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name        string
		input       *user.User
		mockFn      func(mockTx *database.MockTxManager, mockToken *jwttoken.MockJWTToken, mockRepo *userrepository.MockUserRepository, input *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "success",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(mockTx *database.MockTxManager, mockToken *jwttoken.MockJWTToken, mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockUser := &user.User{Email: "test1@mail.com", Password: "$2y$10$WsTQ3C0XLFoAWJNA3kY0AOOkSzZwXF20KVRtjSR18FkS5d20OYwp2"}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), input.Email).Return(mockUser, nil).Times(1)

				mockTx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				mockToken.EXPECT().GenerateAccessToken(mockUser).Return("mock-access-token", nil).Times(1)
				mockToken.EXPECT().GenerateRefreshToken(mockUser).Return("mock-refresh-token", nil).Times(1)

				mockRepo.EXPECT().InsertRefreshTokenTx(gomock.Any(), nil, gomock.Any()).Return(nil).Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "fail invalid email",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(mockTx *database.MockTxManager, mockToken *jwttoken.MockJWTToken, mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), input.Email).Return(nil, errs.ErrInvalidEmailOrPassword).Times(1)
			},
			expectedErr: errs.ErrInvalidEmailOrPassword,
		},
		{
			name:  "fail invalid password",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(mockTx *database.MockTxManager, mockToken *jwttoken.MockJWTToken, mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockUser := &user.User{Email: "test1@mail.com", Password: "invalid_password"}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), input.Email).Return(mockUser, nil).Times(1)
			},
			expectedErr: errs.ErrInvalidEmailOrPassword,
		},
		{
			name:  "fail generate token",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(mockTx *database.MockTxManager, mockToken *jwttoken.MockJWTToken, mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockUser := &user.User{Email: "test1@mail.com", Password: "$2y$10$WsTQ3C0XLFoAWJNA3kY0AOOkSzZwXF20KVRtjSR18FkS5d20OYwp2"}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), input.Email).Return(mockUser, nil).Times(1)

				mockTx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				mockToken.EXPECT().GenerateAccessToken(mockUser).Return("mock-access-token", nil).Times(1)
				mockToken.EXPECT().GenerateRefreshToken(mockUser).Return("", ErrDB).Times(1)
			},
			expectedErr: ErrDB,
		},
		{
			name:  "success",
			input: &user.User{Email: "test1@mail.com", Password: "test_password"},
			mockFn: func(mockTx *database.MockTxManager, mockToken *jwttoken.MockJWTToken, mockRepo *userrepository.MockUserRepository, input *user.User) {
				mockUser := &user.User{Email: "test1@mail.com", Password: "$2y$10$WsTQ3C0XLFoAWJNA3kY0AOOkSzZwXF20KVRtjSR18FkS5d20OYwp2"}
				mockRepo.EXPECT().FindUserByEmail(gomock.Any(), input.Email).Return(mockUser, nil).Times(1)

				mockTx.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(tx *sql.Tx) error) error {
						return fn(nil)
					},
				).Times(1)

				mockToken.EXPECT().GenerateAccessToken(mockUser).Return("mock-access-token", nil).Times(1)
				mockToken.EXPECT().GenerateRefreshToken(mockUser).Return("mock-refresh-token", nil).Times(1)

				mockRepo.EXPECT().InsertRefreshTokenTx(gomock.Any(), nil, gomock.Any()).Return(ErrDB).Times(1)
			},
			expectedErr: ErrDB,
		},
	}

	for _, tc := range testCases {
		mockToken, mockTx, mockRepo, service := setup(t)

		tc.mockFn(mockTx, mockToken, mockRepo, tc.input)

		resp, err := service.Login(context.Background(), tc.input.Email, tc.input.Password)

		if tc.expectedErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, resp)
		}
	}
}

func setup(t *testing.T) (*jwttoken.MockJWTToken, *database.MockTxManager, *userrepository.MockUserRepository, userservice.UserService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockToken := jwttoken.NewMockJWTToken(ctrl)
	mockTx := database.NewMockTxManager(ctrl)
	mockRepo := userrepository.NewMockUserRepository(ctrl)

	service := userservice.NewUserService(mockTx, mockToken, mockRepo)

	return mockToken, mockTx, mockRepo, service
}
