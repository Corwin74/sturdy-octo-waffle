package user

import (
	"context"
	"errors"
	"shop/internal/conf"
	"shop/internal/models"
	"shop/internal/common"
	user_repo_mock "shop/mocks/repository/user"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := user_repo_mock.NewMockUserRepo(ctrl)
	testConfig := &conf.Secrets{
		JwtKey: "test-secret-key",
	}

	usecase := NewUsecase(mockRepo, testConfig)
	ctx := context.Background()

	tests := []struct {
		name          string
		username      string
		password      string
		mockBehavior  func()
		expectedError bool
	}{
		{
			name:     "Success Login Existing User",
			username: "existing_user",
			password: "correct_password",
			mockBehavior: func() {
				hashedPassword, _ := HashPassword("correct_password")
				mockRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(models.User{
						ID:       uuid.New(),
						Name:     "existing_user",
						Password: hashedPassword,
						Balance:  1000,
					}, nil)
			},
			expectedError: false,
		},
		{
			name:     "Wrong Password",
			username: "existing_user",
			password: "wrong_password",
			mockBehavior: func() {
				hashedPassword, _ := HashPassword("correct_password")
				mockRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(models.User{
						ID:       uuid.New(),
						Name:     "existing_user",
						Password: hashedPassword,
						Balance:  1000,
					}, nil)
			},
			expectedError: true,
		},
		{
			name:     "Create New User",
			username: "new_user",
			password: "new_password",
			mockBehavior: func() {
				mockRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(models.User{}, common.ErrNotFound)

				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.New(), nil)
			},
			expectedError: false,
		},
		{
			name:     "Repository Error",
			username: "error_user",
			password: "error_password",
			mockBehavior: func() {
				mockRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(models.User{}, errors.New("repository error"))
			},
			expectedError: true,
		},
		{
			name:     "Hashing Error During User Creation",
			username: "new_user",
			password: string(make([]byte, 100000)), // Очень длинный пароль вызовет ошибку хеширования
			mockBehavior: func() {
				mockRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(models.User{}, common.ErrNotFound)
			},
			expectedError: true,
		},
		{
			name:     "Create User Error",
			username: "new_user",
			password: "new_password",
			mockBehavior: func() {
				mockRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(models.User{}, common.ErrNotFound)
		
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.UUID{}, errors.New("create error"))
			},
			expectedError: true,
		},
		// {
		// 	name:     "Token Generation Error",
		// 	username: "existing_user",
		// 	password: "correct_password",
		// 	mockBehavior: func() {
		// 		hashedPassword, _ := HashPassword("correct_password")
		// 		mockRepo.EXPECT().
		// 			Get(gomock.Any(), gomock.Any()).
		// 			Return(models.User{
		// 				ID:       uuid.New(),
		// 				Name:     "existing_user",
		// 				Password: hashedPassword,
		// 				Balance:  1000,
		// 			}, nil)
		// 		testConfig.JwtKey = string(make([]byte, 1<<31)) // Слишком длинный ключ
		// 	},
		// 	expectedError: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			token, err := usecase.Auth(ctx, tt.username, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "test_password"

	hash1, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash1)

	// Проверяем, что два хеша одного и того же пароля различаются
	hash2, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash2)
	assert.NotEqual(t, hash1, hash2)

	// Проверяем валидацию пароля
	assert.True(t, checkPassword(password, hash1))
	assert.True(t, checkPassword(password, hash2))
	assert.False(t, checkPassword("wrong_password", hash1))

	// Проверяем ошибку при слишком длинном пароле
	password = string(make([]byte, 100))
	hash2, err = HashPassword(password)
	assert.Error(t, err)
}

func TestGenerateTokenForUser(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	token := generateTokenForUser(userID, secret)
	assert.NotEmpty(t, token)

	// Генерируем второй токен для того же пользователя
	token2 := generateTokenForUser(userID, secret)
	assert.NotEmpty(t, token2)
	// Токены должны отличаться даже для одного и того же пользователя
	assert.NotEqual(t, token, token2)
}
