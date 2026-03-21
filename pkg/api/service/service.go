package service

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("field 'repeat' is empty")
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
	//Do it
	// case "m":
	//
	// case "w":
	//
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
	}

	return date.Format(DateFormat), nil

}

func afterNow(date time.Time, now time.Time) bool {
	return date.After(now)
}
