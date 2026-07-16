package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const tokenLifetime = 8 * time.Hour

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type tokenPayload struct {
	PasswordHash string `json:"password_hash"`
	ExpiresAt    int64  `json:"exp"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, fmt.Errorf("метод %s не поддерживается", r.Method))
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeError(w, errors.New("пароль на сервере не установлен"))
		return
	}

	var request signInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, fmt.Errorf("ошибка десериализации JSON: %w", err))
		return
	}

	if subtle.ConstantTimeCompare(
		[]byte(request.Password),
		[]byte(password),
	) != 1 {
		writeError(w, errors.New("неверный пароль"))
		return
	}

	token, err := createToken(password)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, signInResponse{
		Token: token,
	})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

		// Если пароль не задан, аутентификация отключена.
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(
				w,
				"Authentication required",
				http.StatusUnauthorized,
			)
			return
		}

		if err := validateToken(cookie.Value, password); err != nil {
			http.Error(
				w,
				"Authentication required",
				http.StatusUnauthorized,
			)
			return
		}

		next(w, r)
	}
}

func createToken(password string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload := tokenPayload{
		PasswordHash: hashPassword(password),
		ExpiresAt:    time.Now().Add(tokenLifetime).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal JWT header: %w", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal JWT payload: %w", err)
	}

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadPart := base64.RawURLEncoding.EncodeToString(payloadJSON)

	unsignedToken := headerPart + "." + payloadPart
	signature := signJWT(unsignedToken, password)

	return unsignedToken + "." + signature, nil
}

func validateToken(token, password string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("некорректный формат JWT")
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSignature := signJWT(unsignedToken, password)

	if !hmac.Equal(
		[]byte(parts[2]),
		[]byte(expectedSignature),
	) {
		return errors.New("неверная подпись JWT")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("decode JWT payload: %w", err)
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return fmt.Errorf("unmarshal JWT payload: %w", err)
	}

	if payload.ExpiresAt <= time.Now().Unix() {
		return errors.New("срок действия JWT истёк")
	}

	expectedHash := hashPassword(password)

	if subtle.ConstantTimeCompare(
		[]byte(payload.PasswordHash),
		[]byte(expectedHash),
	) != 1 {
		return errors.New("JWT создан для другого пароля")
	}

	return nil
}

func signJWT(unsignedToken, password string) string {
	mac := hmac.New(sha256.New, []byte(password))
	_, _ = mac.Write([]byte(unsignedToken))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
