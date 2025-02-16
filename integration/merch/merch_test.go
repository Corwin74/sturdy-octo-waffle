package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Coins     int `json:"coins"`
	Inventory []struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	} `json:"inventory"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

const (
	baseURL  = "http://localhost:8080"
	username = "testuser"
	password = "testpass"
)

func TestMerchPurchase(t *testing.T) {
	// Создаем HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Получаем токен авторизации
	token := getAuthToken(t, client)
	require.NotEmpty(t, token, "Токен авторизации не должен быть пустым")

	// Проверяем начальное состояние баланса и инвентаря
	initialInfo := getInfo(t, client, token)
	initialCoins := initialInfo.Coins

	// Тестовые сценарии покупки мерча
	testCases := []struct {
		name          string
		itemToBuy     string
		itemPrice     int
		expectedError bool
	}{
		{
			name:           "Успешная покупка футболки",
			itemToBuy:     "t-shirt",
			itemPrice:     80,
			expectedError: false,
		},
		{
			name:           "Успешная покупка кружки",
			itemToBuy:     "cup",
			itemPrice:     20,
			expectedError: false,
		},
		{
			name:           "Покупка несуществующего товара",
			itemToBuy:     "nonexistent-item",
			itemPrice:     0,
			expectedError: true,
		},
	}
	
	var totalSpent int // Добавляем счетчик общих расходов

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Выполняем запрос на покупку
			resp, err := buyItem(client, token, tc.itemToBuy)
			require.NoError(t, err, "Ошибка при выполнении запроса")
			defer resp.Body.Close()

			if tc.expectedError {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Ожидался код ошибки 400")
				return
			}

			assert.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался успешный статус код")

			// Увеличиваем счетчик потраченных монет только для успешных покупок
            totalSpent += tc.itemPrice

            // Проверяем обновленное состояние после покупки
            updatedInfo := getInfo(t, client, token)

            // Проверяем, что монеты были списаны корректно с учетом всех предыдущих покупок
            expectedCoins := initialCoins - totalSpent

			assert.Equal(t, expectedCoins, updatedInfo.Coins, "Неверное количество монет после покупки")

			// Проверяем, что товар добавлен в инвентарь
			found := false
			for _, item := range updatedInfo.Inventory {
				if item.Type == tc.itemToBuy {
					found = true
					assert.Greater(t, item.Quantity, 0, "Количество купленного товара должно быть больше 0")
					break
				}
			}
			assert.True(t, found, "Купленный товар не найден в инвентаре")
		})
	}
}

func getAuthToken(t *testing.T, client *http.Client) string {
	authReq := AuthRequest{
		Username: username,
		Password: password,
	}

	resp, err := makeRequest(client, "POST", "/api/auth", authReq, "")
	require.NoError(t, err, "Ошибка при выполнении запроса авторизации")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Неверный статус код при авторизации")

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	require.NoError(t, err, "Ошибка при декодировании ответа авторизации")

	return authResp.Token
}

func getInfo(t *testing.T, client *http.Client, token string) InfoResponse {
	resp, err := makeRequest(client, "GET", "/api/info", nil, token)
	require.NoError(t, err, "Ошибка при получении информации")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Неверный статус код при получении информации")

	var info InfoResponse
	err = json.NewDecoder(resp.Body).Decode(&info)
	require.NoError(t, err, "Ошибка при декодировании информации")

	return info
}

func buyItem(client *http.Client, token, item string) (*http.Response, error) {
	return makeRequest(client, "GET", fmt.Sprintf("/api/buy/%s", item), nil, token)
}

func makeRequest(client *http.Client, method, path string, body interface{}, token string) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, baseURL+path, bytes.NewBuffer(bodyJSON))
	} else {
		req, err = http.NewRequest(method, baseURL+path, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return client.Do(req)
}