package config

import "time"

type Config struct {
	DBPath string `yaml:"dbpath"`
	Port   string `yaml:"port"`
}

type Task struct {
	ID      int64     `json:"id"`
	Date    time.Time `json:"date"`
	Title   string    `json:"title"`
	Comment string    `json:"comment"`
	Repeat  string    `json:"repeat"`
}
