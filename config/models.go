package config

type Config struct {
	DBPath string `yaml:"dbpath"`
	Port   string `yaml:"port"`
}

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
