package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

		writeJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{
				"error": fmt.Sprintf(
					"метод %s не поддерживается",
					r.Method,
				),
			},
		)
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		log.Println("ошибка авторизации: переменная TODO_PASSWORD не установлена")

		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "внутренняя ошибка сервера",
			},
		)
		return
	}

	var request signInRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": fmt.Sprintf(
					"ошибка десериализации JSON: %v",
					err,
				),
			},
		)
		return
	}

	if subtle.ConstantTimeCompare(
		[]byte(request.Password),
		[]byte(password),
	) != 1 {
		writeJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{
				"error": "неверный пароль",
			},
		)
		return
	}

	token, err := createToken(password)
	if err != nil {
		log.Printf("ошибка создания токена: %v", err)

		writeJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "внутренняя ошибка сервера",
			},
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		signInResponse{
			Token: token,
		},
	)
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
			writeJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{
					"error": "требуется авторизация",
				},
			)
			return
		}

		if err := validateToken(cookie.Value, password); err != nil {
			log.Printf("ошибка проверки токена: %v", err)

			writeJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{
					"error": "требуется авторизация",
				},
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

	if _, err := mac.Write([]byte(unsignedToken)); err != nil {
		log.Printf("ошибка вычисления подписи JWT: %v", err)
	}

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))

	return base64.RawURLEncoding.EncodeToString(sum[:])
}
