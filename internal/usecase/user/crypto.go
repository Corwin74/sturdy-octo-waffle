package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10


// HashPassword -- хеширует пароль для безопасного хранения
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func GenerateArgon2Hash(password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }
    
    time := uint32(1)
    memory := uint32(64 * 1024)
    threads := uint8(4)
    keyLen := uint32(32)
    
    hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
    
    // Кодируем результат в base64
    encodedHash := base64.RawStdEncoding.EncodeToString(hash)
    encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
    
    // Возвращаем хеш в формате: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
    return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
        memory, time, threads, encodedSalt, encodedHash), nil
}

func verifyPassword(password, encodedHash string) (match bool, err error) {
    // Разбираем закодированную строку
    // Ожидаемый формат: $argon2id$v=19$m=65536,t=1,p=4$соль$хеш
    vals := strings.Split(encodedHash, "$")
    if len(vals) != 6 {
        return false, errors.New("некорректный формат хеша")
    }

    var version int
    _, err = fmt.Sscanf(vals[2], "v=%d", &version)
    if err != nil {
        return false, err
    }
    if version != argon2.Version {
        return false, errors.New("несовместимая версия argon2")
    }

    // Получаем параметры
    var memory uint32
    var time uint32
    var threads uint8
    _, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
    if err != nil {
        return false, err
    }

    // Декодируем соль и хеш из base64
    salt, err := base64.RawStdEncoding.DecodeString(vals[4])
    if err != nil {
        return false, err
    }

    hash, err := base64.RawStdEncoding.DecodeString(vals[5])
    if err != nil {
        return false, err
    }

    // Вычисляем хеш для предоставленного пароля
    // с теми же параметрами и солью
    hashToCompare := argon2.IDKey(
        []byte(password),
        salt,
        time,
        memory,
        threads,
        uint32(len(hash)))

    // Сравниваем хеши
    return subtle.ConstantTimeCompare(hash, hashToCompare) == 1, nil
}

func generateTokenForUser(userID uuid.UUID, secret string) string {
	now := time.Now()

	claims := jwt.MapClaims{
		"id":  userID.String(),
		"iat": now.Unix(),                     // время создания токена
		"exp": now.Add(24 * time.Hour).Unix(), // срок действия - 24 часа
		"jti": uuid.New().String(),            // уникальный идентификатор токена
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}
