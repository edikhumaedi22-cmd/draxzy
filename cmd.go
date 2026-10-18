package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Rank     int    `json:"rank"`
	MMR      float64 `json:"mmr"`
}

var users = map[string]*User{
	"player1": {ID: "1", Username: "player1", Password: "pass123", Rank: 5, MMR: 1500},
}

var jwtSecret = []byte("hok-secret-key-2024")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	user, exists := users[req.Username]
	if !exists || user.Password != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)

	resp := LoginResponse{Token: tokenString, User: user}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/login", loginHandler)
	fmt.Println("Auth Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
