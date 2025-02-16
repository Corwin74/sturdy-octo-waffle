package sendcoin

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

const (
    baseURL = "http://localhost:8080"
)

type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string `json:"token"`
}

type SendCoinRequest struct {
    ToUser string `json:"toUser"`
    Amount int    `json:"amount"`
}

type InfoResponse struct {
    Coins       int `json:"coins"`
    Inventory   []struct {
        Type     string `json:"type"`
        Quantity int    `json:"quantity"`
    } `json:"inventory"`
    CoinHistory struct {
        Received []struct {
            FromUser string `json:"fromUser"`
            Amount   int    `json:"amount"`
        } `json:"received"`
        Sent []struct {
            ToUser string `json:"toUser"`
            Amount int    `json:"amount"`
        } `json:"sent"`
    } `json:"coinHistory"`
}

func TestCoinTransfer(t *testing.T) {
    // Создаем HTTP-клиент с таймаутом
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Аутентификация первого пользователя
    user1Token := authenticateUser(t, client, "user1", "password1")
    require.NotEmpty(t, user1Token, "Токен первого пользователя не должен быть пустым")

    // Аутентификация второго пользователя
    user2Token := authenticateUser(t, client, "user2", "password2")
    require.NotEmpty(t, user2Token, "Токен второго пользователя не должен быть пустым")

    // Получаем начальный баланс первого пользователя
    initialBalance1 := getUserBalance(t, client, user1Token)

    // Получаем начальный баланс второго пользователя
    initialBalance2 := getUserBalance(t, client, user2Token)

    // Сумма для перевода
    transferAmount := 50

    // Отправляем монеты от первого пользователя второму
    sendCoins(t, client, user1Token, "user2", transferAmount)

    // Проверяем балансы после перевода
    finalBalance1 := getUserBalance(t, client, user1Token)
    finalBalance2 := getUserBalance(t, client, user2Token)

    // Проверяем корректность изменения балансов
    assert.Equal(t, initialBalance1-transferAmount, finalBalance1, "Неверный баланс отправителя после перевода")
    assert.Equal(t, initialBalance2+transferAmount, finalBalance2, "Неверный баланс получателя после перевода")

    // Проверяем историю транзакций первого пользователя
    info1 := getUserInfo(t, client, user1Token)
    require.GreaterOrEqual(t, len(info1.CoinHistory.Sent), 1, "История отправленных монет должна содержать минимум одну запись")
    assert.Equal(t, "user2", info1.CoinHistory.Sent[len(info1.CoinHistory.Sent)-1].ToUser)
    assert.Equal(t, transferAmount, info1.CoinHistory.Sent[len(info1.CoinHistory.Sent)-1].Amount)

    // Проверяем историю транзакций второго пользователя
    info2 := getUserInfo(t, client, user2Token)
    require.GreaterOrEqual(t, len(info2.CoinHistory.Received), 1, "История полученных монет должна содержать минимум одну запись")
    assert.Equal(t, "user1", info2.CoinHistory.Received[len(info2.CoinHistory.Received)-1].FromUser)
    assert.Equal(t, transferAmount, info2.CoinHistory.Received[len(info2.CoinHistory.Received)-1].Amount)
}

func authenticateUser(t *testing.T, client *http.Client, username, password string) string {
    authReq := AuthRequest{
        Username: username,
        Password: password,
    }

    authReqBody, err := json.Marshal(authReq)
    require.NoError(t, err, "Ошибка при маршалинге запроса аутентификации")

    resp, err := client.Post(
        fmt.Sprintf("%s/api/auth", baseURL),
        "application/json",
        bytes.NewBuffer(authReqBody),
    )
    require.NoError(t, err, "Ошибка при выполнении запроса аутентификации")
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode, "Неверный код ответа при аутентификации")

    var authResp AuthResponse
    err = json.NewDecoder(resp.Body).Decode(&authResp)
    require.NoError(t, err, "Ошибка при декодировании ответа аутентификации")

    return authResp.Token
}

func getUserBalance(t *testing.T, client *http.Client, token string) int {
    info := getUserInfo(t, client, token)
    return info.Coins
}

func getUserInfo(t *testing.T, client *http.Client, token string) InfoResponse {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/info", baseURL), nil)
    require.NoError(t, err, "Ошибка при создании запроса информации")

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

    resp, err := client.Do(req)
    require.NoError(t, err, "Ошибка при выполнении запроса информации")
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode, "Неверный код ответа при получении информации")

    var info InfoResponse
    err = json.NewDecoder(resp.Body).Decode(&info)
    require.NoError(t, err, "Ошибка при декодировании ответа с информацией")

    return info
}

func sendCoins(t *testing.T, client *http.Client, token string, toUser string, amount int) {
    sendReq := SendCoinRequest{
        ToUser: toUser,
        Amount: amount,
    }

    sendReqBody, err := json.Marshal(sendReq)
    require.NoError(t, err, "Ошибка при маршалинге запроса отправки монет")

    req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sendCoin", baseURL), bytes.NewBuffer(sendReqBody))
    require.NoError(t, err, "Ошибка при создании запроса отправки монет")

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    require.NoError(t, err, "Ошибка при выполнении запроса отправки монет")
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode, "Неверный код ответа при отправке монет")
}