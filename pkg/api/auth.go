package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignInRequest struct {
	Password string `json:"password"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	pswrd := todoPassword
	if pswrd != req.Password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}
	res, err := tokenHelper(pswrd)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": res})
}

func hashHelper(s string) string {
	data := []byte(s)

	hash := sha256.Sum256(data)

	encoded := hex.EncodeToString(hash[:])
	return encoded
}

func tokenHelper(password string) (string, error) {
	claims := jwt.MapClaims{
		"hash": hashHelper(password),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(password))
}

func tokenValidate(tokenString string, password string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(password), nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	value, ok := claims["hash"]
	if !ok {
		return false, nil
	}
	hashStr, ok := value.(string)
	if !ok {
		return false, nil
	}
	if hashStr != hashHelper(password) {
		return false, nil
	}
	return true, nil
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(todoPassword) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			var valid bool
			if jwt != "" {
				ok, err := tokenValidate(jwt, todoPassword)
				if err == nil && ok {
					valid = true
				}
			}

			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
