package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type Ranking struct {
	Username string `json:"username"`
	ClassID  int    `json:"c_class"`
	Score    int    `json:"score"`
}
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TwoFacData struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}
type Session struct {
	SessionID       string                 `json:"session_id"`
	SessionMetadata map[string]interface{} `json:"session_metadata"`
	ExpiryDatetime  time.Time              `json:"expiry_datetime"`
}

var sessionStore = struct {
	sync.RWMutex
	sessions map[string]Session
}{sessions: make(map[string]Session)}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=postgres password=123 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func startSession(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		res.WriteHeader(http.StatusOK)
		return
	}

	var username string
	err := json.NewDecoder(req.Body).Decode(&username)
	if err != nil {
		http.Error(res, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	sessionMetadata := map[string]interface{}{
		"username":   username,
		"role":       "user",
		"ip_address": req.RemoteAddr,
		"user_agent": req.UserAgent(),
	}

	sessionID := uuid.NewString()
	expiryDatetime := time.Now().Add(30 * time.Minute)

	session := Session{
		SessionID:       sessionID,
		SessionMetadata: sessionMetadata,
		ExpiryDatetime:  expiryDatetime,
	}

	sessionStore.Lock()
	sessionStore.sessions[sessionID] = session
	sessionStore.Unlock()

	http.SetCookie(res, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: expiryDatetime,
	})

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(session)
}

func getSession(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		res.WriteHeader(http.StatusOK)
		return
	}

	cookie, err := req.Cookie("session_id")
	if err != nil {
		http.Error(res, "Session not found", http.StatusUnauthorized)
		return
	}

	sessionStore.RLock()
	session, exists := sessionStore.sessions[cookie.Value]
	sessionStore.RUnlock()

	if !exists || session.ExpiryDatetime.Before(time.Now()) {
		http.Error(res, "Invalid or expired session", http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(session)
}

func checkSessionExpiry() {
	for {
		time.Sleep(1 * time.Minute)
		sessionStore.Lock()
		for id, session := range sessionStore.sessions {
			if session.ExpiryDatetime.Before(time.Now()) {
				delete(sessionStore.sessions, id)
				fmt.Println("Session expired:", id)
			}
		}
		sessionStore.Unlock()
	}
}

func checkSession(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
	if req.Method == http.MethodOptions {
		res.WriteHeader(http.StatusOK)
		return
	}

	cookie, err := req.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(res, "No session found", http.StatusUnauthorized)
			return
		}
		http.Error(res, "Error retrieving cookie", http.StatusUnauthorized)
		return
	}

	sessionStore.Lock()
	session, sessionExists := sessionStore.sessions[cookie.Value]
	sessionStore.Unlock()

	if !sessionExists || session.ExpiryDatetime.Before(time.Now()) {
		http.Error(res, "Session expired or invalid", http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Session is still valid",
		"session": session.SessionMetadata,
	}
	json.NewEncoder(res).Encode(response)
}

func saltHashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.WriteHeader(http.StatusOK)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	res.Header().Set("Access-Control-Allow-Credentials", "true")

	var loginData LoginData
	err := json.NewDecoder(req.Body).Decode(&loginData)
	if err != nil {
		http.Error(res, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	query := "SELECT encrypted_password FROM accounts WHERE username = $1"
	err = db.QueryRow(query, loginData.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(res, "Invalid credentials", http.StatusUnauthorized)
			log.Println("No rows found for username:", loginData.Username)
			return
		}
		http.Error(res, "Server error", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}

	if !verifyPassword(loginData.Password, hashedPassword) {
		http.Error(res, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.NewString()
	expiryDatetime := time.Now().Add(30 * time.Minute)
	sessionMetadata := map[string]interface{}{
		"username":   loginData.Username,
		"role":       "user", // Example: You can adjust based on your role logic
		"ip_address": req.RemoteAddr,
		"user_agent": req.UserAgent(),
	}

	session := Session{
		SessionID:       sessionID,
		SessionMetadata: sessionMetadata,
		ExpiryDatetime:  expiryDatetime,
	}

	sessionStore.Lock()
	sessionStore.sessions[sessionID] = session
	sessionStore.Unlock()

	http.SetCookie(res, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiryDatetime,
		HttpOnly: true,
		Secure:   false,
	})

	res.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Login successful",
		"session": session,
	}

	json.NewEncoder(res).Encode(response)
}

func logout(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.WriteHeader(http.StatusOK)
		return
	}

	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	res.Header().Set("Access-Control-Allow-Credentials", "true")

	cookie, err := req.Cookie("session_id")
	if err != nil {
		http.Error(res, "Session not found", http.StatusUnauthorized)
		return
	}

	sessionStore.Lock()
	delete(sessionStore.sessions, cookie.Value)
	sessionStore.Unlock()

	http.SetCookie(res, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
	})

	res.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Logged out successfully",
	}
	json.NewEncoder(res).Encode(response)
}

func generateTwoFacSecret(issuer string, email string) (string, string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
	})
	if err != nil {
		log.Fatalf("Error generating 2FA secret: %v", err)
	}

	var imgBase64 string
	if img, err := key.Image(500, 500); err == nil {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			log.Println("Error encoding QR code image:", err)
			return "", "", "", err
		}
		imgBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	}

	return key.Secret(), key.URL(), imgBase64, nil
}

func ValidateTwoFacCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func setupTwoFac(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var twoFacData TwoFacData
	err := json.NewDecoder(req.Body).Decode(&twoFacData)
	if err != nil {
		http.Error(res, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var email string
	var secretkey_2fa string
	query := "SELECT email, secretkey_2fa FROM accounts WHERE username = $1"
	err = db.QueryRow(query, twoFacData.Username).Scan(&email, &secretkey_2fa)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(res, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(res, "Server error", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}

	secret, url, base64img, err := generateTwoFacSecret(email, secretkey_2fa)
	if err != nil {
		log.Println("Error Generate 2FA:", err)
		return
	}
	fmt.Println(url)
	fmt.Println(base64img)

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, `{"message": "2FA setup successful", "secret": "%s"}`, secret)
}

func ValidateTwoFac(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var twoFacData TwoFacData
	err := json.NewDecoder(req.Body).Decode(&twoFacData)
	if err != nil {
		http.Error(res, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var secretkey_2fa string
	query := "SELECT secretkey_2fa FROM accounts WHERE username = $1"
	err = db.QueryRow(query, twoFacData.Username).Scan(&secretkey_2fa)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(res, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(res, "Server error", http.StatusInternalServerError)
		log.Println("Database query error:", err)
		return
	}

	// if !ValidateTwoFacCode(secretkey_2fa, twoFacData.Code) {
	// 	http.Error(res, "Invalid credentials", http.StatusUnauthorized)
	// 	return
	// }

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("2FA validated successfully"))
}

func getRankings(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		res.WriteHeader(http.StatusOK)
		return
	}

	query := "SELECT a.username, c.class_id, s.reward_score FROM scores s INNER JOIN characters c ON s.char_id = c.char_id INNER JOIN accounts a ON c.acc_id = a.acc_id "
	args := []interface{}{}
	reqQuery := req.URL.Query()
	// classIDParam := reqQuery.Get("class_id")

	// Username
	usernameParam := reqQuery.Get("username")
	if usernameParam != "" {
		query += "WHERE a.username = $1 "
		args = append(args, usernameParam)
	}

	// Count
	countParam := reqQuery.Get("count")
	if countParam == "" {
		query += "ORDER BY s.reward_score DESC LIMIT 10 "
	} else {
		query += "ORDER BY s.reward_score DESC LIMIT "
		if usernameParam == "" {
			query += "$1 "
		} else {
			query += "$2 "
		}
		limit, err := strconv.Atoi(countParam)
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return
		}
		args = append(args, limit)
	}

	// Pagination
	pageParam := reqQuery.Get("page")
	itemsPerPage, err := strconv.Atoi(countParam)
	if err != nil {
		http.Error(res, "Invalid itemPerPage parameter", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}
	offset := (page - 1) * itemsPerPage
	if page <= 0 {
		query += "OFFSET 0 "
	} else if page > 0 {
		if usernameParam == "" {
			query += "OFFSET $2 "
		} else {
			query += "OFFSET $3 "
		}
	}
	args = append(args, offset)

	// fmt.Println("Query:", query)
	// fmt.Println("Arguments:", args)
	// print("\n")

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rankings []Ranking
	for rows.Next() {
		var rank Ranking
		if err := rows.Scan(&rank.Username, &rank.ClassID, &rank.Score); err != nil {
			log.Fatal(err)
		}
		rankings = append(rankings, rank)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

	if err := json.NewEncoder(res).Encode(rankings); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// --Generate Hash Password--
	// result, err := saltHashPassword("123!@#")
	// if err != nil {
	// 	fmt.Println("Error hashing password:", err)
	// 	return
	// }
	// fmt.Println("Hashed Password:", result)

	//--Generate 2FA--
	// result_1, result_2, result_3, err := generate2FASecret("WiraApp", "amanda04@example.net")
	// if err != nil {
	// 	fmt.Println("Error hashing password:", err)
	// 	return
	// }
	// fmt.Println("2FA:", result_1)
	// fmt.Println("2FA:", result_2)
	// fmt.Println("2FA:", result_3)

	//Periodically Check Session
	go checkSessionExpiry()

	//API
	http.HandleFunc("/api/start-session", startSession)
	http.HandleFunc("/api/get-session", getSession)
	http.HandleFunc("/api/check-session", checkSession)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/rankings", getRankings)
	http.HandleFunc("/api/login", Login)
	http.HandleFunc("/api/setup-2fa", setupTwoFac)
	http.HandleFunc("/api/validate-2fa", ValidateTwoFac)
	log.Fatal(http.ListenAndServe(":8181", nil))
}
