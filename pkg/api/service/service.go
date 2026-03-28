package service

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/robberhood/final_project/config"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// var day [32]bool
	// var month [13]bool

	if repeat == "" {
		return "", errors.New("field 'repeat' has invalid format")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts_repeat := strings.Split(repeat, " ")
	if parts_repeat[0] != "d" && parts_repeat[0] != "y" && parts_repeat[0] != "m" && parts_repeat[0] != "w" {
		return "", errors.New("field 'repeat' has invalid format")
	}

	switch parts_repeat[0] {
	case "y":
		if len(parts_repeat) != 1 {
			return "", errors.New("field 'repeat' has invalid format")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "d":
		if len(parts_repeat) != 2 {
			return "", errors.New("field 'repeat' has invalid format")
		}
		interval, err := strconv.Atoi(parts_repeat[1])
		if err != nil {
			return "", errors.New("field 'repeat' has invalid format")
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("field 'repeat' has invalid interval")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
	// case "w":
	// 	if len(parts_repeat)!= 2 {
	// 		return "", errors.New("field 'repeat' has invalid format")
	// 	}
	// 	if len(parts_repeat[1])> 7 {
	// 		return "", errors.New("field 'repeat' has invalid format")
	// 	}
	// 	for _, c := range parts_repeat[1] {
	// 		if c < '1' || c > '7' {
	// 			return "", errors.New("field 'repeat' has invalid format")
	// 		}

	// 	}

	// 	return "", errors.New("field 'repeat' has invalid format")

	default:
		return "", errors.New("field 'repeat' has invalid format")
	}

	return date.Format(DateFormat), nil

}

func afterNow(date time.Time, now time.Time) bool {
	return date.After(now)
}

func CheckDate(task *config.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if afterNow(today, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func GenerateToken(secret string) (string, error) {
	claims := jwt.MapClaims{
		"pass_hash": sha256.Sum256([]byte(secret)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				WriteError(w, errors.New("Authentication required"))
				return
			}

			jwtT := cookie.Value

			token, err := jwt.Parse(jwtT, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenMalformed
				}
				return []byte(pass), nil
			})

			if err != nil || !token.Valid {
				WriteError(w, errors.New("Authentication required"))
				return
			}
		}

		next(w, r)
	})
}

func WriteJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
