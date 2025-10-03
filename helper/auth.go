package helper

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func JwtGenerate(userID uint) (string, error) {
	secret := []byte(ENV("JWT_SECRET"))
	expHours, _ := strconv.Atoi(ENV("JWT_EXPIRE_HOUR"))
	if expHours == 0 {
		expHours = 72
	}
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Duration(expHours) * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func JwtParse(tokenString string) (jwt.MapClaims, error) {
	secret := []byte(ENV("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func SendEmail(to, subject, body string) error {
	host := ENV("SMTP_HOST")
	port := ENV("SMTP_PORT")
	user := ENV("SMTP_USERNAME")
	pass := ENV("SMTP_PASSWORD")

	if host == "" || user == "" {
		// SMTP not configured — in prod you must configure. For dev just log or return nil.
		fmt.Printf("SendEmail stub -> to:%s subject:%s body:%s\n", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)
	msg := "From: " + user + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	return smtp.SendMail(addr, auth, user, []string{to}, []byte(msg))
}
